package service

import (
	"github.com/gugabfigueiredo/dream-mail-go/log"
	"github.com/gugabfigueiredo/dream-mail-go/models"
	"testing"
)

type MailTestCase struct {
	name     string
	provider IProvider
	mail     models.Mail
}

func TestMailgunProvider_SendMail(t *testing.T) {

	tests := []MailTestCase{
		{
			name: "testing",
		},
	}

	logger := &log.Logger{}

	for _, tt := range tests {
		tt.provider = mp
	}
}
