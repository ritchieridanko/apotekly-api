package mailer

import (
	"log"
	"sync"

	"github.com/ritchieridanko/apotekly-api/auth/config"
	"gopkg.in/gomail.v2"
)

type Mailer interface {
	Send(message *gomail.Message) (err error)
	Close() (err error)
}

type mailer struct {
	dialer *gomail.Dialer
	sender gomail.SendCloser
	mu     sync.Mutex
}

func NewMailer() Mailer {
	return &mailer{
		dialer: gomail.NewDialer(
			config.MailerGetHost(),
			config.MailerGetPort(),
			config.MailerGetUser(),
			config.MailerGetPass(),
		),
	}
}

// sends email safely with reconnect + retry logic
func (m *mailer) Send(message *gomail.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.ensureConnection(); err != nil {
		return err
	}

	if err := gomail.Send(m.sender, message); err != nil {
		log.Println("WARNING -> failed to send email, retrying:", err.Error())

		// reset connection
		_ = m.sender.Close()
		m.sender = nil

		// reconnect + retry once
		if err := m.ensureConnection(); err != nil {
			return err
		}

		return gomail.Send(m.sender, message)
	}

	return nil
}

// shuts down connection
func (m *mailer) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.sender != nil {
		return m.sender.Close()
	}

	return nil
}

// ensures connection is alive, else reconnects
func (m *mailer) ensureConnection() error {
	if m.sender == nil {
		sender, err := m.dialer.Dial()
		if err != nil {
			return err
		}

		m.sender = sender
		log.Println("SUCCESS -> connected to SMTP server")
	}
	return nil
}
