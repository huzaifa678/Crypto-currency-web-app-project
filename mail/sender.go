package mail

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)


type Sender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string) error
}

type GmailSender struct {
	name string
	senderEmail string
	senderPassword string
}

func NewGmailSender(name, email, password string) *GmailSender {
	return &GmailSender{
		name: name,
		senderEmail: email,
		senderPassword: password,
	}
}

func (g *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string) error {
		e := gomail.NewMessage()

		e.SetHeader("From", g.name + " <" + g.senderEmail + ">")

		e.SetHeader("To", to...)

		if len(cc) > 0 {
			e.SetHeader("Cc", cc...)

		}
		if len(bcc) > 0 {
			e.SetHeader("Bcc", bcc...)
		}
		e.SetHeader("Subject", subject)

		e.SetBody("text/HTML", content)

		for _, file := range attachFiles {
			e.Attach(file)
		}
		d := gomail.NewDialer("smtp.gmail.com", 587, g.senderEmail, g.senderPassword)

		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		if err := d.DialAndSend(e); err != nil {
			return err
		}
		
		return nil
	}

