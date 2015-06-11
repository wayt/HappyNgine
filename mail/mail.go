package mail

import (
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
	return provider.Send(m)
}
