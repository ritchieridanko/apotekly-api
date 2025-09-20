package email

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/internal/ce"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infras/mailer"
	"github.com/ritchieridanko/apotekly-api/auth/internal/utils"
	"go.opentelemetry.io/otel"
	"gopkg.in/gomail.v2"
)

const emailErrorTracer string = "service.email"

type EmailService interface {
	SendPasswordResetToken(ctx context.Context, email, token string) (err error)
	SendVerificationToken(ctx context.Context, email, token string) (err error)
	SendWelcomeMessage(ctx context.Context, email, token string) (err error)
}

type emailService struct {
	from     string
	sender   mailer.Mailer
	template *template.Template
}

func NewService(sender mailer.Mailer) EmailService {
	// get the directory of the current file (email.go)
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)

	tmplDir := filepath.Join(baseDir, "templates", "*.html")
	templates, err := template.ParseGlob(tmplDir)
	if err != nil {
		log.Fatalln("FATAL -> failed to parse html template files:", err)
	}

	return &emailService{config.MailerGetUser(), sender, templates}
}

func (es *emailService) SendPasswordResetToken(ctx context.Context, email, token string) error {
	ctx, span := otel.Tracer(emailErrorTracer).Start(ctx, "SendPasswordResetToken")
	defer span.End()

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
		return ce.NewError(span, ce.CodeEmailTemplateParsing, ce.MsgInternalServer, err)
	}

	message := es.buildMessage([]string{email}, "Password Reset Request", body.String(), "")
	return es.sendEmail(ctx, message)
}

func (es *emailService) SendVerificationToken(ctx context.Context, email, token string) error {
	ctx, span := otel.Tracer(emailErrorTracer).Start(ctx, "SendVerificationToken")
	defer span.End()

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
		return ce.NewError(span, ce.CodeEmailTemplateParsing, ce.MsgInternalServer, err)
	}

	message := es.buildMessage([]string{email}, "Verify Your Email Address", body.String(), "")
	return es.sendEmail(ctx, message)
}

func (es *emailService) SendWelcomeMessage(ctx context.Context, email, token string) error {
	ctx, span := otel.Tracer(emailErrorTracer).Start(ctx, "SendWelcomeMessage")
	defer span.End()

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
		return ce.NewError(span, ce.CodeEmailTemplateParsing, ce.MsgInternalServer, err)
	}

	message := es.buildMessage([]string{email}, "Welcome to Apotekly", body.String(), "")
	return es.sendEmail(ctx, message)
}

func (es *emailService) sendEmail(ctx context.Context, message *gomail.Message) (err error) {
	_, span := otel.Tracer(emailErrorTracer).Start(ctx, "sendEmail")
	defer span.End()

	if err := es.sender.Send(message); err != nil {
		return ce.NewError(span, ce.CodeEmailDelivery, ce.MsgInternalServer, err)
	}
	return nil
}

func (es *emailService) buildMessage(recipients []string, subject, htmlBody, textBody string) (message *gomail.Message) {
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
