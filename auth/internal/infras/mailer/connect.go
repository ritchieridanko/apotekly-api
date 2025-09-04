package mailer

import (
	"log"
	"sync"

	"github.com/ritchieridanko/apotekly-api/auth/config"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer
	sender gomail.SendCloser
	mu     sync.Mutex
}

func NewMailer() *Mailer {
	return &Mailer{
		dialer: gomail.NewDialer(
			config.GetMailerHost(),
			config.GetMailerPort(),
			config.GetMailerUser(),
			config.GetMailerPass(),
		),
	}
}

// Ensures the connection is alive, or else reconnects
func (m *Mailer) ensureConnection() (err error) {
	if m.sender == nil {
		sender, err := m.dialer.Dial()
		if err != nil {
			return err
		}

		m.sender = sender
		log.Println("SUCCESS: connected to SMTP server")
	}
	return nil
}

// Sends an email safely with reconnect + retry logic
func (m *Mailer) Send(message *gomail.Message) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := m.ensureConnection(); err != nil {
		return err
	}

	if err := gomail.Send(m.sender, message); err != nil {
		log.Println("WARNING: send failed, retrying:", err)

		// Reset connection
		_ = m.sender.Close()
		m.sender = nil

		// Reconnect + retry once
		if err := m.ensureConnection(); err != nil {
			return err
		}

		return gomail.Send(m.sender, message)
	}

	return nil
}

// Shuts down the connection cleanly
func (m *Mailer) Close() (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.sender != nil {
		return m.sender.Close()
	}

	return nil
}
