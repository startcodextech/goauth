package account

type (
	UserDto struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Lastname     string `json:"lastname"`
		Email        string `json:"email"`
		PasswordHash string `json:"password"`
		PhotoUrl     string `json:"photo_url,omitempty"`
		Verified     bool   `json:"verified"`
		FacebookID   string `json:"facebook_id,omitempty"`
		GoogleID     string `json:"google_id,omitempty"`
		MicrosoftID  string `json:"microsoft_id,omitempty"`
		AppleID      string `json:"apple_id,omitempty"`
	}

	UserRegisterDto struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Lastname    string `json:"lastname"`
		Email       string `json:"email"`
		Password    string `json:"password"` // Contraseña en texto plano
		FacebookID  string `json:"facebook_id,omitempty"`
		GoogleID    string `json:"google_id,omitempty"`
		MicrosoftID string `json:"microsoft_id,omitempty"`
		AppleID     string `json:"apple_id,omitempty"`
	}

	UserCreateDto struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Lastname     string `json:"lastname"`
		Email        string `json:"email"`
		PasswordHash string `json:"password_hash"`
		FacebookID   string `json:"facebook_id,omitempty"`
		GoogleID     string `json:"google_id,omitempty"`
		MicrosoftID  string `json:"microsoft_id,omitempty"`
		AppleID      string `json:"apple_id,omitempty"`
	}
)
