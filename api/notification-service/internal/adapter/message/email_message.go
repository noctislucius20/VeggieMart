package message

import (
	"crypto/tls"
	"notification-service/config"
	"notification-service/utils/logger"

	mail "github.com/go-gomail/gomail"
)

type EmailMessageInterface interface {
	SendEmailNotification(to, subject, body string) error
}

type emailAttribute struct {
	Username string
	Password string
	Host     string
	Port     int
	Sending  string
	IsTLS    bool
}

// SendEmailNotification implements [EmailMessageInterface].
func (e *emailAttribute) SendEmailNotification(to string, subject string, body string) error {
	logger := logger.NewLogger().Logger()

	mailMessage := mail.NewMessage(func(m *mail.Message) {
		m.SetHeader("From", e.Sending)
		m.SetHeader("To", to)
		m.SetHeader("Subject", subject)

		m.SetBody("text/html", body)
	})

	mailDialer := mail.NewDialer(e.Host, e.Port, e.Username, e.Password)
	mailDialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: e.IsTLS,
	}

	if err := mailDialer.DialAndSend(mailMessage); err != nil {
		logger.Errorf("[EmailMessage-1] SendEmailNotification: %v", err.Error())
		return err
	}

	return nil
}

func NewEmailMessage(cfg *config.Config) EmailMessageInterface {
	return &emailAttribute{
		Username: cfg.EmailConfig.Username,
		Password: cfg.EmailConfig.Password,
		Host:     cfg.EmailConfig.Host,
		Port:     cfg.EmailConfig.Port,
		Sending:  cfg.EmailConfig.Sending,
		IsTLS:    cfg.EmailConfig.IsTLS,
	}
}
