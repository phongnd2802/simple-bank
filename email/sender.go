package email

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

const (
	SMTP_HOST = "smtp.gmail.com"
	SMTP_PORT = 587
)

type gmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

// SendEmail implements EmailSender.
func (sender *gmailSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, file := range attachFiles {
		_, err := e.AttachFile(file)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", file, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, SMTP_HOST)
	return e.Send("smtp.gmail.com:587", smtpAuth)
}

func NewGmailSender(name, email, password string) EmailSender {
	return &gmailSender{
		name:              name,
		fromEmailAddress:  email,
		fromEmailPassword: password,
	}
}
