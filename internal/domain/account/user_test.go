package account

import (
	"github.com/google/uuid"
	"github.com/modernice/goes/test"
	"testing"
)

func TestNew(t *testing.T) {
	test.NewAggregate(t, UserNew, UserAggregate)
}

func TestUser_Create(t *testing.T) {
	user := UserNew(uuid.New())

	userNew := UserCreateDto{
		Name:     "Julio",
		Lastname: "Caicedo",
		Email:    "email@test.mail",
		Password: "@Master.123",
	}

	if err := user.Create(userNew); err != nil {
		t.Fatalf("Create(%q) failed with %q", userNew, err)
	}

	changes := user.AggregateChanges()
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}

	change := changes[0]
	if change.Name() != EventUserCreated {
		t.Fatalf("expected event %q, got %q", EventUserCreated, change.Name())
	}

	createdData, ok := change.Data().(UserCreated)
	if !ok {
		t.Fatalf("expected data of type UserCreated, got %T", change.Data())
	}

	if createdData.Name != userNew.Name || createdData.Lastname != userNew.Lastname || createdData.Email != userNew.Email {
		t.Fatalf("event data does not match: %#v", createdData)
	}

	if createdData.PasswordHash == userNew.Password {
		t.Fatal("password was not hashed")
	}
}
