package account

import (
	"encoding/json"
	"errors"
	"regexp"
)

var (
	nameRegexp  = regexp.MustCompile(NameOrLastnameRegexp)
	emailRegexp = regexp.MustCompile(EmailRegexp)
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
// This function is used to create a new User in the system.
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

// ID function returns the id of the User.
func (u User) ID() string { return u.id }

// Name function returns the name of the User.
func (u User) Name() string { return u.name }

// Lastname function returns the lastname of the User.
func (u User) Lastname() string { return u.lastname }

// Email function returns the email of the User.
func (u User) Email() string { return u.email }

// PasswordHash function returns the passwordHash of the User.
func (u User) PasswordHash() string { return u.passwordHash }

// PhotoUrl function returns the photoUrl of the User.
func (u User) PhotoUrl() string { return u.photoUrl }

// Verified function returns the verified of the User.
func (u User) Verified() bool { return u.verified }

// FacebookID function returns the facebookId of the User.
func (u User) FacebookID() string { return u.facebookId }

// GoogleID function returns the googleId of the User.
func (u User) GoogleID() string { return u.googleId }

// MicrosoftID function returns the microsoftId of the User.
func (u User) MicrosoftID() string { return u.microsoftId }

// AppleID function returns the appleId of the User.
func (u User) AppleID() string { return u.appleId }

// Marshal marshals the User into a slice of bytes.
func (u User) Marshal() []byte {
	jsonMap := map[string]interface{}{
		"id":            u.id,
		"name":          u.name,
		"lastname":      u.lastname,
		"email":         u.email,
		"photo_url":     u.photoUrl,
		"verified":      u.verified,
		"facebook_id":   u.facebookId,
		"google_id":     u.googleId,
		"microsoft_id":  u.microsoftId,
		"apple_id":      u.appleId,
		"password_hash": u.passwordHash,
	}

	raw, _ := json.Marshal(jsonMap)
	return raw
}

// UnmarshalFromMap unmarshals the User from a map.
func (u User) UnmarshalFromMap(data map[string]interface{}) error {
	if id, ok := data["id"].(string); ok {
		u.id = id
	}
	if name, ok := data["name"].(string); ok {
		u.name = name
	}
	if lastname, ok := data["lastname"].(string); ok {
		u.lastname = lastname
	}
	if email, ok := data["email"].(string); ok {
		u.email = email
	}
	if photoUrl, ok := data["photo_url"].(string); ok {
		u.photoUrl = photoUrl
	}
	if verified, ok := data["verified"].(bool); ok {
		u.verified = verified
	}
	if facebookId, ok := data["facebook_id"].(string); ok {
		u.facebookId = facebookId
	}
	if googleId, ok := data["google_id"].(string); ok {
		u.googleId = googleId
	}
	if microsoftId, ok := data["microsoft_id"].(string); ok {
		u.microsoftId = microsoftId
	}
	if appleId, ok := data["apple_id"].(string); ok {
		u.appleId = appleId
	}
	return nil
}

// validateName function validates a name.
func (u User) validateName(name string) error {
	if !nameRegexp.MatchString(name) {
		return errors.New("name is not valid")
	}
	return nil
}

// validateLastname function validates a lastname.
func (u User) validateLastname(lastname string) error {
	if err := u.validateName(lastname); err != nil {
		return errors.New("lastname is not valid")
	}
	return nil
}

// validateEmail function validates an email.
func (u User) validateEmail(email string) error {
	if !emailRegexp.MatchString(email) {
		return errors.New("email is not valid")
	}
	return nil
}
