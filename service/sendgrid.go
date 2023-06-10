package service

import (
	"github.com/gugabfigueiredo/dream-mail-go/log"
	"github.com/sendgrid/sendgrid-go"
	sgMail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridProvider struct {
	Logger *log.Logger
	Client *sendgrid.Client
}

func newSendgridProvider(apiKey string, logger *log.Logger) *SendgridProvider {
	client := sendgrid.NewSendClient(apiKey)

	return &SendgridProvider{
		Logger: logger,
		Client: client,
	}
}

func (s *SendgridProvider) SendMail() error {

	logger := s.Logger.C()

	from := sgMail.NewEmail("Tira", "tiramisu@example.com") // Change to your verified sender
	subject := "Sending with Twilio SendGrid is Fun"
	to := sgMail.NewEmail("Sou", "souffle@example.com") // Change to your recipient
	plainTextContent := "and easy to do anywhere, even with Go"
	message := sgMail.NewSingleEmailPlainText(from, subject, to, plainTextContent)

	resp, err := s.Client.Send(message)
	if err != nil {
		logger.E("failed to send message", "err", err.Error())
	}

	logger.I("email request successful", "status", resp.StatusCode, "headers", resp.Headers)

	return nil
}
