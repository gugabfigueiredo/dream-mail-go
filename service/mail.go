package service

import (
	sp "github.com/SparkPost/gosparkpost"
	"github.com/gugabfigueiredo/dream-mail-go/log"
	"github.com/gugabfigueiredo/dream-mail-go/models"
)

type IService interface {
	QueueMail(mail *models.Mail)
}

type IProvider interface {
	SendMail(mail *models.Mail) error
}

type Config struct {
	SMTPConfig
	SESConfig
	SendgridConfig
	sp.Config
}

type Service struct {
	Logger       *log.Logger
	Providers    []IProvider
	mailingQueue chan *models.Mail
}

var done chan bool

func NewService(providers []IProvider, logger *log.Logger) *Service {

	s := &Service{
		Logger:       logger,
		Providers:    providers,
		mailingQueue: make(chan *models.Mail, 100),
	}

	done = make(chan bool)
	go s.sendQueued()

	return s
}

func (s *Service) QueueMail(mail *models.Mail) {
	s.mailingQueue <- mail
}

func (s *Service) sendQueued() {
	for {
		select {
		case mail := <-s.mailingQueue:
			var err error
			logger := s.Logger.C("mailID", mail.ID, "from", mail.From.Addr, "to", mail.To)
			for _, provider := range s.Providers {
				err = provider.SendMail(mail)
				if err == nil {
					break
				}
				logger.E("unable to send to provider", "err", err)
			}
		case <-done:
			close(s.mailingQueue)
			return
		}
	}
}

func (s *Service) Quit() {
	done <- true
	close(done)
}
