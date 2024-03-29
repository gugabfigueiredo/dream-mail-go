package service

import (
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/gugabfigueiredo/dream-mail-go/models"
	log "github.com/gugabfigueiredo/tiny-go-log"

	//go get -u github.com/aws/aws-sdk-go
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type SESConfig struct {
	Region string `json:"region"`
}

type SESProvider struct {
	Logger  *log.Logger
	Client  sesiface.SESAPI
	Session *session.Session
}

func NewSESProvider(cfg SESConfig, logger *log.Logger) *SESProvider {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region)},
	)

	if err != nil {
		logger.F("Error creating AWS session: %s", err)
	}

	// Create an SES session.
	svc := ses.New(sess)

	return &SESProvider{
		Logger:  logger.C("provider", "ses"),
		Client:  svc,
		Session: sess,
	}
}

func (s *SESProvider) SendMail(mail *models.Mail) error {

	logger := s.Logger.C("from", mail.From.Addr, "to", mail.To, "subject", mail.Subject, "id", mail.ID)

	sesInput := s.buildInput(mail)

	// Attempt to send the email.
	result, err := s.Client.SendRawEmail(sesInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			logger.E("unable to send email", "err", aerr.Error(), "code", aerr.Code())
		} else {
			logger.E("unknown error", "err", err.Error())
		}

		return err
	}

	logger.I("SES SendRawEmail successful", "result", result)

	return nil
}

func (s *SESProvider) buildInput(mail *models.Mail) *ses.SendRawEmailInput {

	var tos []string
	for _, to := range mail.To {
		tos = append(tos, to.Addr)
	}

	msg := buildSMTPMessage(mail, tos)

	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: msg,
		},
	}

	return input
}
