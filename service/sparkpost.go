package service

import (
	sp "github.com/SparkPost/gosparkpost"
	"github.com/gugabfigueiredo/dream-mail-go/models"
	log "github.com/gugabfigueiredo/tiny-go-log"
	"github.com/pkg/errors"
)

type ISPClient interface {
	Send(*sp.Transmission) (id string, res *sp.Response, err error)
}

type SparkpostProvider struct {
	Logger *log.Logger
	Client ISPClient
}

func NewSparkpostProvider(cfg sp.Config, logger *log.Logger) *SparkpostProvider {

	var client sp.Client
	err := client.Init(&cfg)
	if err != nil {
		logger.F("SparkPost client init failed: %s\n", err)
	}

	return &SparkpostProvider{
		Logger: logger,
		Client: &client,
	}
}

func (s *SparkpostProvider) SendMail(mail *models.Mail) error {

	logger := s.Logger.C("from", mail.From.Addr, "to", mail.To, "subject", mail.Subject, "id", mail.ID)

	// Create a Transmission
	tx := s.buildTransmission(mail)

	id, resp, err := s.Client.Send(tx)
	if err != nil {
		logger.E("unable to send email", "err", err, "statusCode", resp.HTTP.StatusCode)
		return errors.New("unable to send email")
	}

	logger.I("sparkpost transmission sent", "id", id, "statusCode", resp.HTTP.StatusCode)

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
			HTML:        mail.HTML,
			Attachments: attachments,
		},
	}

	return tx
}
