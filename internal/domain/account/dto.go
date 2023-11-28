package account

type (
	UserCreateDto struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Lastname     string `json:"lastname"`
		Email        string `json:"email"`
		PasswordHash string `json:"password"`
		FacebookID   string `json:"facebook_id,omitempty"`
		GoogleID     string `json:"google_id,omitempty"`
		MicrosoftID  string `json:"microsoft_id,omitempty"`
		AppleID      string `json:"apple_id,omitempty"`
	}
)
