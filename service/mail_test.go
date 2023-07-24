package service

import (
	"github.com/gugabfigueiredo/dream-mail-go/models"
	log "github.com/gugabfigueiredo/tiny-go-log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockProvider struct {
	CallCount  int
	CalledWith []*models.Mail

	CallsBeforeError int
	Error            error

	ExpectedCallCount  int
	ExpectedCalledWith []*models.Mail
}

func (m *MockProvider) SendMail(mail *models.Mail) error {
	m.CallCount++
	m.CalledWith = append(m.CalledWith, mail)
	if m.CallsBeforeError != 0 && m.CallCount > m.CallsBeforeError {
		return m.Error
	}
	return nil
}

func TestService_QueueAndSendMail(t *testing.T) {
	tests := []struct {
		name      string
		mails     []*models.Mail
		providers []IProvider
	}{
		{
			name: "queue and send mail with success - should call first provider only",
			mails: []*models.Mail{
				{
					From: models.Email{
						Addr: "sender@domain.com",
					},
					To: []models.Email{
						{
							Addr: "recipient@domain.com",
						},
					},
					Subject: "Test email 1",
					Text:    "This is a test email 1",
				},
			},
			providers: []IProvider{
				&MockProvider{
					ExpectedCallCount: 1,
					ExpectedCalledWith: []*models.Mail{
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 1",
							Text:    "This is a test email 1",
						},
					},
				},
				&MockProvider{
					ExpectedCallCount:  0,
					ExpectedCalledWith: nil,
				},
			},
		},
		{
			name: "queue and send mail with fail-over - should call next provider",
			mails: []*models.Mail{
				{
					From: models.Email{
						Addr: "sender@domain.com",
					},
					To: []models.Email{
						{
							Addr: "recipient@domain.com",
						},
					},
					Subject: "Test email 1",
					Text:    "This is a test email 1",
				},
			},
			providers: []IProvider{
				&MockProvider{
					CallsBeforeError:  -1,
					Error:             errors.New("error sending email"),
					ExpectedCallCount: 1,
					ExpectedCalledWith: []*models.Mail{
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 1",
							Text:    "This is a test email 1",
						},
					},
				},
				&MockProvider{
					ExpectedCallCount: 1,
					ExpectedCalledWith: []*models.Mail{
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 1",
							Text:    "This is a test email 1",
						},
					},
				},
			},
		},
		{
			name: "queue and send multiple mail - call provider with success - should call provider many times",
			mails: []*models.Mail{
				{
					From: models.Email{
						Addr: "sender@domain.com",
					},
					To: []models.Email{
						{
							Addr: "recipient@domain.com",
						},
					},
					Subject: "Test email 1",
					Text:    "This is a test email 1",
				},
				{
					From: models.Email{
						Addr: "sender@domain.com",
					},
					To: []models.Email{
						{
							Addr: "recipient@domain.com",
						},
					},
					Subject: "Test email 2",
					Text:    "This is a test email 2",
				},
			},
			providers: []IProvider{
				&MockProvider{
					ExpectedCallCount: 2,
					ExpectedCalledWith: []*models.Mail{
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 1",
							Text:    "This is a test email 1",
						},
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 2",
							Text:    "This is a test email 2",
						},
					},
				},
				&MockProvider{
					ExpectedCallCount:  0,
					ExpectedCalledWith: nil,
				},
			},
		},
		{
			name: "queue and send multiple mail - call providers with errors - should call different providers",
			mails: []*models.Mail{
				{
					From: models.Email{
						Addr: "sender@domain.com",
					},
					To: []models.Email{
						{
							Addr: "recipient@domain.com",
						},
					},
					Subject: "Test email 1",
					Text:    "This is a test email 1",
				},
				{
					From: models.Email{
						Addr: "sender@domain.com",
					},
					To: []models.Email{
						{
							Addr: "recipient@domain.com",
						},
					},
					Subject: "Test email 2",
					Text:    "This is a test email 2",
				},
				{
					From: models.Email{
						Addr: "sender@domain.com",
					},
					To: []models.Email{
						{
							Addr: "recipient@domain.com",
						},
					},
					Subject: "Test email 3",
					Text:    "This is a test email 3",
				},
				{
					From: models.Email{
						Addr: "sender@domain.com",
					},
					To: []models.Email{
						{
							Addr: "recipient@domain.com",
						},
					},
					Subject: "Test email 4",
					Text:    "This is a test email 4",
				},
				{
					From: models.Email{
						Addr: "sender@domain.com",
					},
					To: []models.Email{
						{
							Addr: "recipient@domain.com",
						},
					},
					Subject: "Test email 5",
					Text:    "This is a test email 5",
				},
			},
			providers: []IProvider{
				&MockProvider{
					CallsBeforeError:  2,
					Error:             errors.New("error sending email"),
					ExpectedCallCount: 5,
					ExpectedCalledWith: []*models.Mail{
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 1",
							Text:    "This is a test email 1",
						},
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 2",
							Text:    "This is a test email 2",
						},
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 3",
							Text:    "This is a test email 3",
						},
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 4",
							Text:    "This is a test email 4",
						},
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 5",
							Text:    "This is a test email 5",
						},
					},
				},
				&MockProvider{
					CallsBeforeError:  1,
					Error:             errors.New("error sending email"),
					ExpectedCallCount: 3,
					ExpectedCalledWith: []*models.Mail{
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 3",
							Text:    "This is a test email 3",
						},
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 4",
							Text:    "This is a test email 4",
						},
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 5",
							Text:    "This is a test email 5",
						},
					},
				},
				&MockProvider{
					ExpectedCallCount: 2,
					ExpectedCalledWith: []*models.Mail{
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 4",
							Text:    "This is a test email 4",
						},
						{
							From: models.Email{
								Addr: "sender@domain.com",
							},
							To: []models.Email{
								{
									Addr: "recipient@domain.com",
								},
							},
							Subject: "Test email 5",
							Text:    "This is a test email 5",
						},
					},
				},
			},
		},
	}

	logger := log.New(&log.Config{
		Context:               "dmail-go",
		ConsoleLoggingEnabled: false,
		EncodeLogsAsJson:      true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService(tt.providers, logger)

			for _, mail := range tt.mails {
				s.QueueMail(mail)
			}

			time.Sleep(1 * time.Second)
			s.Quit()

			for _, provider := range tt.providers {
				p := provider.(*MockProvider)
				assert.Equal(t, p.ExpectedCalledWith, p.CalledWith, "provider CalledWIth does not match expected")
			}
		})
	}
}
