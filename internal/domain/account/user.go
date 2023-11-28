package account

import (
	"errors"
	"regexp"
)

type (
	// User represents a person who uses the application.
	User struct {
		id           string
		name         string
		lastname     string
		email        string
		passwordHash string
		photoUrl     string
		verified     bool
		facebookId   string
		googleId     string
		microsoftId  string
		appleId      string
	}
)

// NewUser creates a new User.
// This function is used to create a new User instance.
func NewUser() User {
	return User{}
}

// Create function validates and creates a new User.
func (u *User) Create(user UserCreateDto) (err error) {
	if user.ID == "" {
		return errors.New("id is required")
	}

	if err = u.validateName(user.Name); err != nil {
		return
	}

	if err = u.validateLastname(user.Lastname); err != nil {
		return
	}

	if err = u.validateEmail(user.Email); err != nil {
		return
	}

	u.id = user.ID
	u.name = user.Name
	u.lastname = user.Lastname
	u.email = user.Email
	u.passwordHash = user.PasswordHash
	u.verified = false

	return nil
}

func (u *User) ID() string           { return u.id }
func (u *User) Name() string         { return u.name }
func (u *User) Lastname() string     { return u.lastname }
func (u *User) Email() string        { return u.email }
func (u *User) PasswordHash() string { return u.passwordHash }
func (u *User) PhotoUrl() string     { return u.photoUrl }
func (u *User) Verified() bool       { return u.verified }
func (u *User) FacebookID() string   { return u.facebookId }
func (u *User) GoogleID() string     { return u.googleId }
func (u *User) MicrosoftID() string  { return u.microsoftId }
func (u *User) AppleID() string      { return u.appleId }

// SetID function sets the id of the User.
func (u *User) SetID(id string) *User {
	u.id = id
	return u
}

// SetName function sets the name of the User.
func (u *User) SetName(name string) *User {
	u.name = name
	return u
}

// SetLastname function sets the lastname of the User.
func (u *User) SetLastname(lastname string) *User {
	u.lastname = lastname
	return u
}

// SetEmail function sets the email of the User.
func (u *User) SetEmail(email string) *User {
	u.email = email
	return u
}

// SetPasswordHash function sets the passwordHash of the User.
func (u *User) SetPasswordHash(passwordHash string) *User {
	u.passwordHash = passwordHash
	return u
}

// SetPhotoUrl function sets the photoUrl of the User.
func (u *User) SetPhotoUrl(photoUrl string) *User {
	u.photoUrl = photoUrl
	return u
}

// SetVerified function sets the verified of the User.
func (u *User) SetVerified(verified bool) *User {
	u.verified = verified
	return u
}

// SetFacebookID function sets the facebookId of the User.
func (u *User) SetFacebookID(facebookId string) *User {
	u.facebookId = facebookId
	return u
}

// SetGoogleID function sets the googleId of the User.
func (u *User) SetGoogleID(googleId string) *User {
	u.googleId = googleId
	return u
}

// SetMicrosoftID function sets the microsoftId of the User.
func (u *User) SetMicrosoftID(microsoftId string) *User {
	u.microsoftId = microsoftId
	return u
}

// SetAppleID function sets the appleId of the User.
func (u *User) SetAppleID(appleId string) *User {
	u.appleId = appleId
	return u
}

// validateName function validates a name.
func (u *User) validateName(name string) error {
	if ok, err := regexp.MatchString(NameOrLastnameRegexp, name); !ok || err != nil {
		return errors.New("name is not valid")
	}

	return nil
}

// validateLastname function validates a lastname.
func (u *User) validateLastname(lastname string) error {
	if err := u.validateName(lastname); err != nil {
		return errors.New("lastname is not valid")
	}
	return nil
}

// validateEmail function validates an email.
func (u *User) validateEmail(email string) error {
	if ok, err := regexp.MatchString(EmailRegexp, email); !ok || err != nil {
		return errors.New("email is not valid")
	}

	return nil
}
