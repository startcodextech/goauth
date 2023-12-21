package brevo

type (
	ErrorModel struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	PostSendFailed struct {
		ErrorModel
		UnexistingEmails  []string `json:"unexistingEmails,omitempty"`
		WithoutListEmails []string `json:"withoutListEmails,omitempty"`
		BlackListedEmails []string `json:"blackListedEmails,omitempty"`
	}

	CreateSmtpEmail struct {
		MessageId  string   `json:"messageId"`
		MessageIds []string `json:"messageIds"`
	}
)
