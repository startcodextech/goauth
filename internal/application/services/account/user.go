package account

import (
	"context"
	"errors"
	"github.com/google/uuid"
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
	// UserService represents the server for managing users.
	UserService struct {
		repository account.UserRepository
	}
)

// Ensure UserService implements the account.UserService interface.
var _ account.UserService = (*UserService)(nil)

// NewUserService creates a new user server.
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
func (s UserService) isValidPassword(pwd string) bool {
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
func (s UserService) Create(ctx context.Context, data account.UserRegisterDto) (string, error) {

	emailRegistred, err := s.IsExists(ctx, data.Email)
	if err != nil {
		return "", err
	}
	if emailRegistred {
		return "", errors.New("user already exists")
	}

	if !s.isValidPassword(data.Password) {
		return "", errors.New("password is not valid")
	}

	pwdBytes, err := bcrypt.GenerateFromPassword([]byte(data.Password), passwordCost)
	if err != nil {
		return "", err
	}

	user := account.NewUser()

	id := uuid.New().String()

	err = user.Create(account.UserCreateDto{
		ID:           id,
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
		return "", err
	}

	s.repository.Save(ctx, user)
	err = s.repository.Commit(ctx)
	if err != nil {
		return "", err
	}

	return id, nil
}

// IsExists checks if a user exists.
// It returns true if the user exists, otherwise it returns false.
// It returns an error if the operation fails.
// The email parameter is the email of the user to check.
func (s UserService) IsExists(ctx context.Context, email string) (bool, error) {
	filter := map[string]interface{}{
		"email": email,
	}

	results, err := s.repository.Find(ctx, filter)
	if err != nil {
		return false, err
	}

	return len(results) > 0, nil
}
