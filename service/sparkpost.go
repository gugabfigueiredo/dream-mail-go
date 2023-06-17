package service

import (
	sp "github.com/SparkPost/gosparkpost"
	"github.com/gugabfigueiredo/dream-mail-go/log"
	"github.com/gugabfigueiredo/dream-mail-go/models"
	"github.com/pkg/errors"
)

type SparkpostProvider struct {
	Logger *log.Logger
	Client *sp.Client
}

func newSparkpostProvider(cfg *sp.Config, logger *log.Logger) *SparkpostProvider {

	var client sp.Client
	err := client.Init(cfg)
	if err != nil {
		logger.F("SparkPost client init failed: %s\n", err)
	}

	return &SparkpostProvider{
		Logger: logger,
		Client: &client,
	}
}

func (s *SparkpostProvider) SendMail(mail *models.Mail) error {

	// Create a Transmission
	tx := s.buildTransmission(mail)

	id, _, err := s.Client.Send(tx)
	if err != nil {
		s.Logger.E("unable to send email", "err", err)
		return errors.New("unable to send email")
	}

	// The second value returned from Send
	// has more info about the HTTP response, in case
	// you'd like to see more than the Transmission id.
	s.Logger.I("transmission sent", "id", id)

	return nil
}

func (s *SparkpostProvider) buildTransmission(mail *models.Mail) *sp.Transmission {

	var recipients []string
	for _, recipient := range mail.To {
		recipients = append(recipients, recipient.Addr)
	}

	var attachments []sp.Attachment
	for _, attachment := range mail.Attachments {
		attachments = append(attachments, sp.Attachment{
			Filename: attachment.Name,
			MIMEType: attachment.Type,
			B64Data:  attachment.Data,
		})
	}

	// Map the Mail struct to SparkPostTransmission struct
	tx := &sp.Transmission{
		Recipients: recipients,
		Content: sp.Content{
			From:        mail.From.Addr,
			Subject:     mail.Subject,
			Text:        mail.Text,
			HTML:        mail.Html,
			Attachments: attachments,
		},
	}

	return tx
}
