package email

import (
	"bytes"
	"time"

	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
)

func (es *emailService) SendPasswordResetToken(email, token string) error {
	tracer := EmailErrorTracer + ": SendPasswordResetToken()"

	data := struct {
		Email string
		URL   string
		Year  int
	}{
		Email: email,
		URL:   utils.GenerateURLWithTokenQuery("/auth/reset-password", token),
		Year:  time.Now().UTC().Year(),
	}

	var body bytes.Buffer
	if err := es.template.ExecuteTemplate(&body, "password_reset.html", data); err != nil {
		return ce.NewError(ce.ErrCodeParsing, ce.ErrMsgInternalServer, tracer, err)
	}

	message := es.BuildMessage([]string{email}, "Password Reset Request", body.String(), "")
	return es.SendEmail(message)
}

func (es *emailService) SendVerificationToken(email, token string) error {
	tracer := EmailErrorTracer + ": SendVerificationToken()"

	data := struct {
		Email string
		URL   string
		Year  int
	}{
		Email: email,
		URL:   utils.GenerateURLWithTokenQuery("/auth/verify-email", token),
		Year:  time.Now().UTC().Year(),
	}

	var body bytes.Buffer
	if err := es.template.ExecuteTemplate(&body, "email_verification.html", data); err != nil {
		return ce.NewError(ce.ErrCodeParsing, ce.ErrMsgInternalServer, tracer, err)
	}

	message := es.BuildMessage([]string{email}, "Verify Your Email Address", body.String(), "")
	return es.SendEmail(message)
}

func (es *emailService) SendWelcomeMessage(email, token string) error {
	tracer := EmailErrorTracer + ": SendWelcomeMessage()"

	data := struct {
		Email string
		URL   string
		Year  int
	}{
		Email: email,
		URL:   utils.GenerateURLWithTokenQuery("/auth/verify-email", token),
		Year:  time.Now().UTC().Year(),
	}

	var body bytes.Buffer
	if err := es.template.ExecuteTemplate(&body, "welcome.html", data); err != nil {
		return ce.NewError(ce.ErrCodeParsing, ce.ErrMsgInternalServer, tracer, err)
	}

	message := es.BuildMessage([]string{email}, "Welcome to Apotekly", body.String(), "")
	return es.SendEmail(message)
}
