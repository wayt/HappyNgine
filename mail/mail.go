package mail

import (
	"errors"
	"github.com/wayt/happyngine/env"
)

var provider EmailProviderInterface

func init() {
	switch env.Get("EMAIL_PROVIDER") {
	case "mailgun":
		provider = newMailgun()
	}
}

type Email struct {
	From        string
	Subject     string
	Content     string
	ContentHTML string
	To          []string
}

type EmailProviderInterface interface {
	Send(Email) error
}

func Send(m Email) error {
	if provider == nil {

		return errors.New("No email provider supplied")
	}
	return provider.Send(m)
}
