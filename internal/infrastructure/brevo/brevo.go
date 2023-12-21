package brevo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"net/http"
	"os"
)

type (
	Brevo struct {
		key    string
		url    string
		logger *zap.Logger
	}

	SendTo struct {
		Email string
		Name  string
	}
)

func New(logger *zap.Logger) Brevo {
	return Brevo{
		key:    os.Getenv("BREVO_API"),
		url:    "https://api.brevo.com/v3",
		logger: logger,
	}
}

func (b Brevo) SendTemplateEmail(ctx context.Context, test bool, templateID int64, emails []SendTo, params map[string]interface{}) error {
	if test {
		url := fmt.Sprintf("%s%s%d%s", b.url, "/smtp/templates/", templateID, "/sendTest")
		agent := fiber.Post(url)
		agent.ContentType("application/json")
		agent.Add("api-key", b.key)
		agent.JSON(fiber.Map{
			"params": params,
		})
		status, body, errs := agent.Bytes()
		if len(errs) > 0 {
			return errs[0]
		}
		if status == http.StatusNoContent {
			return nil
		}
		if status == http.StatusBadRequest {
			var e PostSendFailed
			err := json.Unmarshal(body, &e)
			if err != nil {
				return err
			}
			return errors.New(fmt.Sprintf(
				"%s  UnexistingEmails: %v  WithoutListEmails: %v BlackListedEmails: %v",
				e.Message,
				e.UnexistingEmails,
				e.WithoutListEmails,
				e.BlackListedEmails,
			))
		}
		if status == http.StatusNotFound {
			return errors.New("template ID not found")
		}
		return nil
	}

	agent := fiber.Post(fmt.Sprintf("%s%s", b.url, "/smtp/email"))
	agent.ContentType("application/json")
	agent.Add("api-key", b.key)

	agent.JSON(fiber.Map{
		"to":         emails,
		"params":     params,
		"templateId": templateID,
	})
	status, body, errs := agent.Bytes()
	if len(errs) > 0 {
		return errs[0]
	}

	if status == http.StatusCreated || status == http.StatusAccepted {
		return nil
	}

	var e ErrorModel
	err := json.Unmarshal(body, &e)
	if err != nil {
		return err
	}
	return errors.New(e.Message)
}
