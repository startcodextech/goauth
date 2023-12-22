package brevo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
)

const (
	baseURL      = "https://api.brevo.com/v3"
	emailTestURL = baseURL + "/smtp/templates/%d/sendTest"
	emailURL     = baseURL + "/smtp/email"
	contentType  = "application/json"
)

type (
	Brevo struct {
		key string
	}
)

func New() *Brevo {
	return &Brevo{
		key: os.Getenv("BREVO_API"),
	}
}

func (b *Brevo) SendTemplateEmail(ctx context.Context, test bool, templateID int64, emails []SendTo, params map[string]interface{}) error {
	url := fmt.Sprintf(emailTestURL, templateID)
	if !test {
		url = emailURL
	}

	agent := fiber.Post(url).ContentType(contentType).Add("api-key", b.key)
	if test {
		agent.JSON(
			fiber.Map{
				"params": params,
			},
		)
	} else {
		agent.JSON(
			fiber.Map{
				"to":         emails,
				"params":     params,
				"templateId": templateID,
			},
		)
	}

	status, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return errs[0]
	}

	switch status {
	case http.StatusNoContent, http.StatusCreated, http.StatusAccepted:
		return nil
	case http.StatusBadRequest:
		var e PostSendFailed
		if err := json.Unmarshal(body, &e); err != nil {
			return err
		}
		return errors.New(e.Message)
	case http.StatusNotFound:
		return errors.New("template ID not found")
	default:
		var e ErrorModel
		if err := json.Unmarshal(body, &e); err != nil {
			return err
		}
		return errors.New(e.Message)
	}
}
