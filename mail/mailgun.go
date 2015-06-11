package mail

import (
	"github.com/mailgun/mailgun-go"
	"github.com/wayt/happyngine/env"
)

type mailgunProvider struct {
	mg mailgun.Mailgun
}

func newMailgun() EmailProviderInterface {

	return &mailgunProvider{
		mg: mailgun.NewMailgun(env.Get("MAILGUN_DOMAIN"), env.Get("MAILGUN_APIKEY"), env.Get("MAILGUN_PUBLIC_APIKEY")),
	}

}

func (p *mailgunProvider) Send(e Email) error {

	msg := mailgun.NewMessage(e.From, e.Subject, e.Content, e.To...)

	if len(e.ContentHTML) > 0 {
		msg.SetHtml(e.ContentHTML)
	}

	_, _, err := p.mg.Send(msg)
	return err
}
