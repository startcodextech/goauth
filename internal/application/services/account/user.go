package account

import (
	"context"
	"errors"
	"github.com/startcodextech/goauth/internal/domain/account"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"unicode"
)

var (
	// passwordSpecialCharRegexp is a regular expression that matches special characters in passwords.
	passwordSpecialCharRegexp = regexp.MustCompile(`[@[\]^_!"#$%&'()*+,-./:;{}<>|=~?]`)
)

const (
	// passwordCost is the cost of the bcrypt algorithm.
	passwordCost = 14
)

type (
	// UserService represents the service for managing users.
	UserService struct {
		repository account.UserRepository
	}
)

// Ensure UserService implements the account.UserService interface.
var _ account.UserService = (*UserService)(nil)

// NewUserService creates a new user service.
// It requires a user repository to be injected.
func NewUserService(repository account.UserRepository) UserService {
	return UserService{
		repository: repository,
	}
}

// isValidPassword checks if the password is valid.
// A password is valid if it has at least 8 characters, one uppercase letter,
// one lowercase letter, one number and one special character.
// The special characters are: @[]^_!"#$%&'()*+,-./:;{}<>|=~?
// The special characters are defined in the passwordSpecialCharRegexp variable.
func (UserService) isValidPassword(pwd string) bool {
	if len(pwd) < 8 {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range pwd {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case passwordSpecialCharRegexp.MatchString(string(char)):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// Create creates a new user.
func (s UserService) Create(ctx context.Context, data account.UserRegisterDto) error {

	if !s.isValidPassword(data.Password) {
		return errors.New("password is not valid")
	}

	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(data.Password), passwordCost)
	if err != nil {
		return err
	}

	user := account.NewUser()

	err = user.Create(account.UserCreateDto{
		ID:           data.ID,
		Name:         data.Name,
		Lastname:     data.Lastname,
		Email:        data.Email,
		PasswordHash: string(pwdBytes),
		FacebookID:   data.FacebookID,
		GoogleID:     data.GoogleID,
		MicrosoftID:  data.MicrosoftID,
		AppleID:      data.AppleID,
	})
	if err != nil {
		return err
	}

	s.repository.Save(ctx, user)
	err = s.repository.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
