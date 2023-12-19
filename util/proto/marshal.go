package proto

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"reflect"
)

type (
	// ProtoMarshal is the default Protocol Buffers marshaler.
	ProtoMarshal struct {
		NewUUID      func() string
		GenerateName func(v interface{}) string
	}

	// NoProtoMessageError is returned when the given value does not implement proto.Message.
	NoProtoMessageError struct {
		v interface{}
	}
)

func (e NoProtoMessageError) Error() string {
	rv := reflect.ValueOf(e.v)
	if rv.Kind() != reflect.Ptr {
		return "v is not proto.Message, you must pass pointer value to implement proto.Message"
	}

	return "v is not proto.Message"
}

var _ cqrs.CommandEventMarshaler = (*ProtoMarshal)(nil)

func (m ProtoMarshal) newUUID() string {
	if m.NewUUID != nil {
		return m.NewUUID()
	}

	// default
	return watermill.NewUUID()
}

func (m ProtoMarshal) Name(cmdOrEvent interface{}) string {
	if m.GenerateName != nil {
		return m.GenerateName(cmdOrEvent)
	}

	return cqrs.FullyQualifiedStructName(cmdOrEvent)
}

func (m ProtoMarshal) Marshal(v interface{}) (*message.Message, error) {
	protoMsg, ok := v.(proto.Message)
	if !ok {
		return nil, errors.WithStack(NoProtoMessageError{v})
	}

	b, err := proto.Marshal(protoMsg)
	if err != nil {
		return nil, err
	}

	msg := message.NewMessage(
		m.newUUID(),
		b,
	)
	msg.Metadata.Set("name", m.Name(v))

	return msg, nil
}

func (ProtoMarshal) Unmarshal(msg *message.Message, v interface{}) (err error) {
	return proto.Unmarshal(msg.Payload, v.(proto.Message))
}

func (m ProtoMarshal) NameFromMessage(msg *message.Message) string {
	return msg.Metadata.Get("name")
}
