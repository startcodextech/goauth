package account

type (
	UserCreateDto struct {
		Name     string `json:"name"`
		Lastname string `json:"lastname"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
)
