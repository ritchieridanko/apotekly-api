package email

import (
	"html/template"
	"log"

	"github.com/ritchieridanko/apotekly-api/auth/config"
	"github.com/ritchieridanko/apotekly-api/auth/internal/infras/mailer"
	"github.com/ritchieridanko/apotekly-api/auth/pkg/ce"
)

const EmailErrorTracer = ce.EmailTracer

type EmailService interface {
	SendPasswordResetToken(email, token string) (err error)
}

type emailService struct {
	from     string
	sender   *mailer.Mailer
	template *template.Template
}

func NewEmailService(sender *mailer.Mailer) EmailService {
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalln("FATAL: unable to parse html template files:", err)
	}

	return &emailService{config.GetMailerUser(), sender, templates}
}
