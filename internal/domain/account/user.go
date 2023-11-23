package account

import (
	"errors"
	"github.com/google/uuid"
	"github.com/modernice/goes/aggregate"
	"github.com/modernice/goes/command"
	"github.com/modernice/goes/command/handler"
	"github.com/modernice/goes/event"
	"github.com/startcodextech/goutils/password"
	"github.com/startcodextech/goutils/validator"
	"log"
)

type (
	User struct {
		*aggregate.Base
		*handler.BaseHandler

		name         string
		lastname     string
		email        string
		passwordHash string
		photoUrl     string
	}
)

const (
	UserAggregate = "account.user"
)

func UserNew(id uuid.UUID) *User {
	var user *User
	user = &User{
		Base: aggregate.New(UserAggregate, id),
		BaseHandler: handler.NewBase(
			handler.BeforeHandle[any](func(ctx command.Ctx[any]) error {
				log.Printf("Handling %q command ... [user=%s]", ctx.Name(), id)
				return nil
			}),
			handler.AfterHandle[any](func(c command.Ctx[any]) {
				//
			}),
		),
	}

	// Register events
	event.ApplyWith(user, user.onCreated, EventUserCreated)

	// Register command
	command.ApplyWith(user, user.Create, UserCreatedCmd)

	return user
}

func (u *User) Create(payload UserCreateDto) error {

	if !validator.IsNameOrLastname(payload.Name) {
		return errors.New("the name provided is not valid.")
	}

	if !validator.IsNameOrLastname(payload.Lastname) {
		return errors.New("the last name provided is not valid.")
	}

	if !validator.IsEmail(payload.Email) {
		return errors.New("the provided email does not meet the validation criteria.")
	}

	if !validator.IsValidPassword(payload.Password) {
		return errors.New("the provided password does not meet the validation criteria.")
	}

	user := UserCreated{
		Name:         payload.Name,
		Lastname:     payload.Lastname,
		Email:        payload.Email,
		PasswordHash: password.HashPasswordString(payload.Password, 14),
	}

	aggregate.Next(u, EventUserCreated, user)
	return nil
}

func (u *User) onCreated(event event.Of[UserCreated]) {
	u.name = event.Data().Name
	u.lastname = event.Data().Lastname
	u.email = event.Data().Email
	u.passwordHash = event.Data().PasswordHash
}
