package service

import (
	log "github.com/gugabfigueiredo/tiny-go-log"
	"testing"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
	"github.com/gugabfigueiredo/dream-mail-go/models"
	"github.com/stretchr/testify/assert"
)

type MockSESClient struct {
	sesiface.SESAPI

	Output     *ses.SendRawEmailOutput
	Error      error
	CalledWith *ses.SendRawEmailInput
}

func (m *MockSESClient) SendRawEmail(input *ses.SendRawEmailInput) (*ses.SendRawEmailOutput, error) {
	m.CalledWith = input
	return m.Output, m.Error
}

func TestSESProvider_SendMail(t *testing.T) {

	tests := []struct {
		name                 string
		mail                 *models.Mail
		mockOutput           *ses.SendRawEmailOutput
		mockError            error
		expectedSESInputData string
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
			mockOutput:           &ses.SendRawEmailOutput{},
			mockError:            nil,
			expectedSESInputData: "From: sender@domain.com\r\nTo: recipient@domain.com\r\nSubject: Test Subject\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=\"dmailboundary\"\r\n\r\n--dmailboundary\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\nTest Text\r\n\r\n--dmailboundary\r\n",
			expectedError:        nil,
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
			mockOutput:           &ses.SendRawEmailOutput{},
			mockError:            nil,
			expectedSESInputData: "From: sender@domain.com\r\nTo: recipient@domain.com\r\nSubject: Test Subject\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=\"dmailboundary\"\r\n\r\n--dmailboundary\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n<h1>Hello World!</h1>\r\n\r\n--dmailboundary\r\n",
			expectedError:        nil,
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
			mockOutput:           &ses.SendRawEmailOutput{},
			mockError:            nil,
			expectedSESInputData: "From: sender@domain.com\r\nTo: recipient@domain.com\r\nSubject: Test Subject\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=\"dmailboundary\"\r\n\r\n--dmailboundary\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n<h1>Hello World!</h1>\r\n\r\n--dmailboundary\r\nContent-Type: text/plain; charset=\"utf-8\"\r\nContent-Disposition: attachment;filename=\"test.txt\"\r\n\r\ntest\r\n\r\n--dmailboundary\r\n",
			expectedError:        nil,
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
			mockOutput:           &ses.SendRawEmailOutput{},
			mockError:            nil,
			expectedSESInputData: "From: sender@domain.com\r\nTo: recipient@domain.com\r\nSubject: Test Subject\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=\"dmailboundary\"\r\n\r\n--dmailboundary\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n<h1>Hello World!</h1>\r\n\r\n--dmailboundary\r\nContent-Type: text/plain; charset=\"utf-8\"\r\nContent-Disposition: attachment;filename=\"test.txt\"\r\n\r\ntest\r\n\r\n--dmailboundary\r\n",
			expectedError:        nil,
		},
	}

	logger := log.New(&log.Config{
		Context:               "dmail-go",
		ConsoleLoggingEnabled: false,
		EncodeLogsAsJson:      true,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock SES client
			mockSvc := &MockSESClient{
				Output: tt.mockOutput,
				Error:  tt.mockError,
			}

			// Create a new instance of our service
			provider := &SESProvider{
				Logger: logger,
				Client: mockSvc,
			}

			err := provider.SendMail(tt.mail)

			// Assertions
			assert.Equal(t, tt.expectedSESInputData, string(mockSvc.CalledWith.RawMessage.Data), "Unexpected SES input")
			assert.Equal(t, tt.expectedError, err, "Unexpected error")
		})
	}
}
