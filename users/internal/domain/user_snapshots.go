package domain

type (
	UserV1 struct {
		Email        string
		Phone        string
		PasswordHash string
		Name         string
		LastName     string
		ProfilePhoto string
		Enabled      bool
		Verified     bool
	}
)

func (UserV1) SnapshotName() string { return "account.UserV1" }
