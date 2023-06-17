package service

import (
	"github.com/sendgrid/rest"
	"testing"

	"github.com/gugabfigueiredo/dream-mail-go/log"
	"github.com/gugabfigueiredo/dream-mail-go/models"

	"github.com/sendgrid/sendgrid-go"
	sgMail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"
)

type MockSendGridClient struct {
	Client     *sendgrid.Client
	Response   *rest.Response
	Error      error
	CalledWith *sgMail.SGMailV3
}

func (m *MockSendGridClient) Send(msg *sgMail.SGMailV3) (*rest.Response, error) {
	m.CalledWith = msg
	return m.Response, m.Error
}

func TestSendgridProvider_SendMail(t *testing.T) {
	tests := []struct {
		name                 string
		mail                 *models.Mail
		mockResponse         *rest.Response
		mockError            error
		expectedSendGridData *sgMail.SGMailV3
		expectedError        error
	}{
		{
			name: "send email plaintext",
			mail: &models.Mail{
				ID: "1234",
				From: models.Email{
					Addr: "sender@domain.com",
					Name: "Sender",
				},
				To: []models.Email{
					{
						Addr: "recipient@domain.com",
						Name: "Recipient",
					},
				},
				Subject: "Test Subject",
				Text:    "Test Text",
			},
			mockResponse: &rest.Response{},
			mockError:    nil,
			expectedSendGridData: &sgMail.SGMailV3{
				From: sgMail.NewEmail("Sender", "sender@domain.com"),
				Personalizations: []*sgMail.Personalization{
					{
						To: []*sgMail.Email{
							sgMail.NewEmail("Recipient", "recipient@domain.com"),
						},
						CC:                  make([]*sgMail.Email, 0),
						BCC:                 make([]*sgMail.Email, 0),
						Headers:             make(map[string]string),
						Substitutions:       make(map[string]string),
						CustomArgs:          make(map[string]string),
						DynamicTemplateData: make(map[string]interface{}),
						Categories:          make([]string, 0),
					},
				},
				Subject: "Test Subject",
				Content: []*sgMail.Content{
					{
						Type:  "text/plain",
						Value: "Test Text",
					},
				},
				Attachments: make([]*sgMail.Attachment, 0),
			},
			expectedError: nil,
		},
		{
			name: "send email html",
			mail: &models.Mail{
				ID: "1234",
				From: models.Email{
					Addr: "sender@domain.com",
					Name: "Sender",
				},
				To: []models.Email{
					{
						Addr: "recipient@domain.com",
						Name: "Recipient",
					},
				},
				Subject: "Test Subject",
				HTML:    "<h1>Hello World!</h1>",
			},
			mockResponse: &rest.Response{},
			mockError:    nil,
			expectedSendGridData: &sgMail.SGMailV3{
				From: sgMail.NewEmail("Sender", "sender@domain.com"),
				Personalizations: []*sgMail.Personalization{
					{
						To: []*sgMail.Email{
							sgMail.NewEmail("Recipient", "recipient@domain.com"),
						},
						CC:                  make([]*sgMail.Email, 0),
						BCC:                 make([]*sgMail.Email, 0),
						Headers:             make(map[string]string),
						Substitutions:       make(map[string]string),
						CustomArgs:          make(map[string]string),
						DynamicTemplateData: make(map[string]interface{}),
						Categories:          make([]string, 0),
					},
				},
				Subject: "Test Subject",
				Content: []*sgMail.Content{
					{
						Type:  "text/html",
						Value: "<h1>Hello World!</h1>",
					},
				},
				Attachments: make([]*sgMail.Attachment, 0),
			},
			expectedError: nil,
		},
		{
			name: "send email with attachment",
			mail: &models.Mail{
				ID: "1234",
				From: models.Email{
					Addr: "sender@domain.com",
					Name: "Sender",
				},
				To: []models.Email{
					{
						Addr: "recipient@domain.com",
						Name: "Recipient",
					},
				},
				Subject: "Test Subject",
				HTML:    "<h1>Hello World!</h1>",
				Attachments: []models.Attachment{
					{
						Name: "test.txt",
						Type: "text/plain",
						Data: "test",
					},
				},
			},
			mockResponse: &rest.Response{},
			mockError:    nil,
			expectedSendGridData: &sgMail.SGMailV3{
				From: sgMail.NewEmail("Sender", "sender@domain.com"),
				Personalizations: []*sgMail.Personalization{
					{
						To: []*sgMail.Email{
							sgMail.NewEmail("Recipient", "recipient@domain.com"),
						},
						CC:                  make([]*sgMail.Email, 0),
						BCC:                 make([]*sgMail.Email, 0),
						Headers:             make(map[string]string),
						Substitutions:       make(map[string]string),
						CustomArgs:          make(map[string]string),
						DynamicTemplateData: make(map[string]interface{}),
						Categories:          make([]string, 0),
					},
				},
				Subject: "Test Subject",
				Content: []*sgMail.Content{
					{
						Type:  "text/html",
						Value: "<h1>Hello World!</h1>",
					},
				},
				Attachments: []*sgMail.Attachment{
					{
						Filename:    "test.txt",
						Type:        "text/plain",
						Content:     "test",
						Disposition: "attachment",
					},
				},
			},
			expectedError: nil,
		},
	}

	logger := log.New(&log.Config{
		Context:               "dmail-go",
		ConsoleLoggingEnabled: false,
		EncodeLogsAsJson:      true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock SendGrid client
			mockClient := &MockSendGridClient{
				Response: tt.mockResponse,
				Error:    tt.mockError,
			}

			// Create a new instance of our service
			provider := &SendgridProvider{
				Logger: logger,
				Client: mockClient,
			}

			err := provider.SendMail(tt.mail)

			// Assertions
			assert.Equal(t, tt.expectedSendGridData.From, mockClient.CalledWith.From, "Unexpected Seder")
			assert.Equal(t, tt.expectedSendGridData.Personalizations, mockClient.CalledWith.Personalizations, "Unexpected Recipients")
			assert.Equal(t, tt.expectedSendGridData.Subject, mockClient.CalledWith.Subject, "Unexpected Subject")
			assert.Equal(t, tt.expectedSendGridData.Content, mockClient.CalledWith.Content, "Unexpected Content")
			assert.Equal(t, tt.expectedSendGridData.Attachments, mockClient.CalledWith.Attachments, "Unexpected Attachments")
			assert.Equal(t, tt.expectedError, err, "Unexpected error")
		})
	}
}
