package service

import (
	"fmt"
	"github.com/gugabfigueiredo/dream-mail-go/log"
	"github.com/gugabfigueiredo/dream-mail-go/models"
	"net/smtp"
)

type SMTPConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type SMTPProvider struct {
	Addr   string
	Auth   smtp.Auth
	Logger *log.Logger
}

func NewSMTPProvider(cfg SMTPConfig, logger *log.Logger) *SMTPProvider {
	return &SMTPProvider{
		Addr:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Auth:   smtp.PlainAuth("", cfg.User, cfg.Pass, cfg.Host),
		Logger: logger.C("provider", "smtp"),
	}
}

func (s *SMTPProvider) SendMail(mail *models.Mail) error {

	var tos []string
	for _, to := range mail.To {
		tos = append(tos, to.Addr)
	}

	msg := buildSMTPMessage(mail, tos)

	err := smtp.SendMail(s.Addr, s.Auth, mail.From.Addr, tos, msg)
	if err != nil {
		s.Logger.E("unable to send email", "err", err)
		return err
	}

	s.Logger.I("email sent", "from", mail.From.Addr, "to", tos)
	return nil
}
