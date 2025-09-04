package email

import (
	"html/template"
	"log"
	"path/filepath"
	"runtime"

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
	// Get the directory of the current file (interface.go)
	_, filename, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(filename)

	tmplDir := filepath.Join(baseDir, "templates", "*.html")
	templates, err := template.ParseGlob(tmplDir)
	if err != nil {
		log.Fatalln("FATAL: unable to parse html template files:", err)
	}

	return &emailService{config.GetMailerUser(), sender, templates}
}
