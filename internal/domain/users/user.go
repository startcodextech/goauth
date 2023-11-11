package users

import (
	"fmt"
	"github.com/startcodextech/goerrors"
	"github.com/startcodextech/goevents/ddd"
	"github.com/startcodextech/goevents/eventsourcing"
	"github.com/startcodextech/goutils/validator"
	"regexp"
	"strings"
)

const (
	UserAggregate = "account.User"

	ErrUserVerifiedCreated  = goerrors.Error("User cannot be created if verified")
	ErrUserEmailValid       = goerrors.Error("The email provided is not valid")
	ErrUserPasswordNotValid = goerrors.Error("Password does not meet security criteria")
	ErrUserName             = goerrors.Error("")
	ErrUserLastName         = goerrors.Error("")
)

var (
	//Minimum length: Generally, a strong password is at least 8 characters.
	//Inclusion of uppercase and lowercase letters: It must contain at least one uppercase and one lowercase letter.
	//Numbers: Include at least one number.
	//Special characters: Include at least one special character (such as @, #, $, etc.).
	passwordRegex = regexp.MustCompile("^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[^\\da-zA-Z]).{8,}$")
)

type (
	User struct {
		eventsourcing.Aggregate
		Email        string
		Phone        string
		PasswordHash string
		Name         string
		LastName     string
		ProfilePhoto string
		Enabled      bool
		Verified     bool
	}
)

var _ interface {
	eventsourcing.EventApplier
	eventsourcing.Snapshotter
} = (*User)(nil)

func NewUser(id string) *User {
	return &User{
		Aggregate: eventsourcing.NewAggregate(id, UserAggregate),
	}
}

func (User) Key() string { return UserAggregate }

func (u *User) CreateUser(id, email, phone, password, name, lastname string) (ddd.Event, error) {
	if u.Verified {
		return nil, ErrUserVerifiedCreated
	}

	if emailValid := validator.Email(email).IsValid(); !emailValid {
		return nil, ErrUserEmailValid
	}

	if !passwordRegex.MatchString(password) {
		return nil, ErrUserPasswordNotValid
	}

	if len(strings.TrimSpace(name)) < 4 {
		return nil, ErrUserName
	}
	if len(strings.TrimSpace(lastname)) < 4 {
		return nil, ErrUserLastName
	}

	u.AddEvent(UserCreatedEvent, &UserCreated{
		ID:       id,
		Email:    email,
		Phone:    phone,
		Name:     name,
		LastName: lastname,
	})

	return ddd.NewEvent(UserCreatedEvent, u), nil
}

func (u *User) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *UserCreated:
		u.Email = payload.Email
		u.Phone = payload.Phone
		u.Name = payload.Name
		u.LastName = payload.LastName
	default:
		return fmt.Errorf("%T received the event %s with unexpected payload %T", u, event.EventName(), payload)
	}
	return nil
}

func (u *User) ApplySnapshot(snapshot eventsourcing.Snapshot) error {
	switch ss := snapshot.(type) {
	case *UserV1:
		u.Email = ss.Email
		u.Phone = ss.Phone
		u.PasswordHash = ss.PasswordHash
		u.Name = ss.Name
		u.LastName = ss.LastName
		u.Enabled = ss.Enabled
		u.Verified = ss.Verified
	default:
		return fmt.Errorf("%T received the unexpected snapshot %T", u, snapshot)
	}

	return nil
}

func (u *User) ToSnapshot() eventsourcing.Snapshot {
	return &UserV1{
		Email:        u.Email,
		Phone:        u.Phone,
		PasswordHash: u.PasswordHash,
		Name:         u.Name,
		LastName:     u.LastName,
		ProfilePhoto: u.ProfilePhoto,
		Enabled:      u.Enabled,
		Verified:     u.Verified,
	}
}
