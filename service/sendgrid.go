package service

import (
	"fmt"
	"github.com/gugabfigueiredo/dream-mail-go/models"
	log "github.com/gugabfigueiredo/tiny-go-log"
	"github.com/pkg/errors"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	sgHelper "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridConfig struct {
	ApiKey string `json:"api_key"`
}

type ISendgridClient interface {
	Send(msg *sgHelper.SGMailV3) (*rest.Response, error)
}

type SendgridProvider struct {
	Logger *log.Logger
	Client ISendgridClient
}

func NewSendgridProvider(cfg SendgridConfig, logger *log.Logger) *SendgridProvider {
	client := sendgrid.NewSendClient(cfg.ApiKey)

	return &SendgridProvider{
		Logger: logger.C("provider", "sendgrid"),
		Client: client,
	}
}

func (s *SendgridProvider) SendMail(mail *models.Mail) error {

	logger := s.Logger.C("from", mail.From.Addr, "to", mail.To, "subject", mail.Subject, "id", mail.ID)

	// Create an instance of SGMailV3
	sgMail := s.buildSGMailV3(mail)

	resp, err := s.Client.Send(sgMail)
	if err != nil {
		logger.E("failed to send message", "err", err.Error(), "status", resp.StatusCode)
		return errors.New(fmt.Sprintf("failed to send message: %s", err.Error()))
	}

	switch resp.StatusCode {
	case 200, 201, 202:
		logger.I("sgMail request successful", "status", resp.StatusCode)
	default:
		logger.E("sgMail request failed", "status", resp.StatusCode)
		return errors.New(fmt.Sprintf("sgMail request failed: %d", resp.StatusCode))
	}

	return nil
}

func (s *SendgridProvider) buildSGMailV3(mail *models.Mail) *sgHelper.SGMailV3 {

	sgMail := sgHelper.NewV3Mail()

	// Set the from address
	sgMail.SetFrom(sgHelper.NewEmail(mail.From.Name, mail.From.Addr))

	// Set the subject
	sgMail.Subject = mail.Subject

	// Set the plain text content
	if mail.Text != "" {
		sgMail.AddContent(sgHelper.NewContent("text/plain", mail.Text))
	}

	// Set the HTML content
	if mail.HTML != "" {
		sgMail.AddContent(sgHelper.NewContent("text/html", mail.HTML))
	}

	// Set the recipients
	personalization := sgHelper.NewPersonalization()

	for _, recipient := range mail.To {
		personalization.AddTos(sgHelper.NewEmail(recipient.Name, recipient.Addr))
	}

	sgMail.AddPersonalizations(personalization)

	// Set the attachments
	for _, attachment := range mail.Attachments {
		sgAtt := sgHelper.NewAttachment()
		sgAtt.SetContent(attachment.Data)
		sgAtt.SetType(attachment.Type)
		sgAtt.SetFilename(attachment.Name)
		sgAtt.SetDisposition("attachment")
		sgMail.AddAttachment(sgAtt)
	}

	return sgMail
}
