package persistence

import "context"

type (
	Entity interface {
		Marshal() []byte
	}

	Model interface {
		ID() string
		Marshal() []byte
		UnmarshalFromMap(map[string]interface{}) error
	}

	Repository[T Entity] interface {
		StartTx(ctx context.Context) (interface{}, error)
		SetTx(session interface{})
		Save(ctx context.Context, model Model)
		Find(ctx context.Context, filter map[string]interface{}) ([]Model, error)
		Delete(ctx context.Context, ID string)
		Commit(ctx context.Context) error
	}
)
