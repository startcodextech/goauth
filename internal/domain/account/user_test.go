package account

import "testing"

func TestCreateUser(t *testing.T) {

	dto := UserCreateDto{
		ID:           "1",
		Name:         "John",
		Lastname:     "Doe",
		Email:        "user@test.app",
		PasswordHash: "14$%&/()=?¡",
	}

	user := NewUser()

	if err := user.Create(dto); err != nil {
		t.Error(err)
	}
}

func TestCreateUser_FailInvalidEmail(t *testing.T) {
	dto := UserCreateDto{
		ID:           "1",
		Name:         "John",
		Lastname:     "Doe",
		Email:        "user@testapp",
		PasswordHash: "14$%&/()=?¡",
	}

	user := NewUser()

	err := user.Create(dto)
	if err == nil {
		t.Errorf("Expected an error for invalid email, got none")
	} else {
		expectedErrorMsg := "email is not valid"
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error '%s', got '%s'", expectedErrorMsg, err.Error())
		}
	}
}

func TestCreateUser_FailInvalidLastName(t *testing.T) {
	dto := UserCreateDto{
		ID:           "1",
		Name:         "John",
		Lastname:     "D",
		Email:        "user@test.app",
		PasswordHash: "14$%&/()=?¡",
	}

	user := NewUser()

	err := user.Create(dto)
	if err == nil {
		t.Errorf("Expected an error for invalid lastname, got none")
	} else {
		expectedErrorMsg := "lastname is not valid"
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error '%s', got '%s'", expectedErrorMsg, err.Error())
		}
	}
}

func TestCreateUser_FailInvalidName(t *testing.T) {
	dto := UserCreateDto{
		ID:           "1",
		Name:         "J",
		Lastname:     "Doe",
		Email:        "user@test.app",
		PasswordHash: "14$%&/()=?¡",
	}

	user := NewUser()

	err := user.Create(dto)
	if err == nil {
		t.Errorf("Expected an error for invalid name, got none")
	} else {
		expectedErrorMsg := "name is not valid"
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error '%s', got '%s'", expectedErrorMsg, err.Error())
		}
	}
}
func TestCreateUser_FaildInvalidID(t *testing.T) {
	dto := UserCreateDto{
		Name:         "John",
		Lastname:     "Doe",
		Email:        "user@test.app",
		PasswordHash: "14$%&/()=?¡",
	}

	user := NewUser()

	err := user.Create(dto)
	if err == nil {
		t.Errorf("Expected an error for invalid id, got none")
	} else {
		expectedErrorMsg := "id is required"
		if err.Error() != expectedErrorMsg {
			t.Errorf("Expected error '%s', got '%s'", expectedErrorMsg, err.Error())
		}
	}

}
