package email

import (
	"encoding/base64"

	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
	"gopkg.in/gomail.v2"
)

func (es *emailService) SendEmail(message *gomail.Message) (err error) {
	tracer := EmailErrorTracer + ": SendEmail()"

	if err := es.sender.Send(message); err != nil {
		return ce.NewError(ce.ErrCodeEmail, ce.ErrMsgInternalServer, tracer, err)
	}
	return nil
}

func (es *emailService) BuildMessage(recipients []string, subject, htmlBody, textBody string) (message *gomail.Message) {
	message = gomail.NewMessage()
	message.SetHeader("From", es.from)
	message.SetHeader("To", recipients...)
	message.SetHeader("Subject", encodeSubject(subject))

	// plain-text + HTML fallback
	if textBody == "" {
		textBody = "Please view this email in an HTML-compatible client."
	}

	message.SetBody("text/plain", textBody)
	message.AddAlternative("text/html", htmlBody)

	return message
}

func encodeSubject(subject string) (encodedSubject string) {
	return "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(subject)) + "?="
}
