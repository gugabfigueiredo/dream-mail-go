package service

import (
	"testing"

	"github.com/gugabfigueiredo/dream-mail-go/log"
	"github.com/gugabfigueiredo/dream-mail-go/models"

	sparkpost "github.com/SparkPost/gosparkpost"
	"github.com/stretchr/testify/assert"
)

type MockSparkPostClient struct {
	Client     sparkpost.Client
	Response   *sparkpost.Response
	Error      error
	CalledWith *sparkpost.Transmission
}

func (m *MockSparkPostClient) Send(transmission *sparkpost.Transmission) (string, *sparkpost.Response, error) {
	m.CalledWith = transmission
	return "", m.Response, m.Error
}

func TestSparkPostProvider_SendMail(t *testing.T) {
	tests := []struct {
		name                  string
		mail                  *models.Mail
		mockResponse          *sparkpost.Response
		mockError             error
		expectedSparkPostData *sparkpost.Transmission
		expectedError         error
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
			mockResponse: &sparkpost.Response{},
			mockError:    nil,
			expectedSparkPostData: &sparkpost.Transmission{
				Content: sparkpost.Content{
					Text:        "Test Text",
					HTML:        "",
					From:        "sender@domain.com",
					Subject:     "Test Subject",
					Attachments: nil,
				},
				Recipients: []string{"recipient@domain.com"},
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
			mockResponse: &sparkpost.Response{},
			mockError:    nil,
			expectedSparkPostData: &sparkpost.Transmission{
				Content: sparkpost.Content{
					Text:        "",
					HTML:        "<h1>Hello World!</h1>",
					From:        "sender@domain.com",
					Subject:     "Test Subject",
					Attachments: nil,
				},
				Recipients: []string{"recipient@domain.com"},
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
			mockResponse: &sparkpost.Response{},
			mockError:    nil,
			expectedSparkPostData: &sparkpost.Transmission{
				Content: sparkpost.Content{
					Text:    "",
					HTML:    "<h1>Hello World!</h1>",
					From:    "sender@domain.com",
					Subject: "Test Subject",
					Attachments: []sparkpost.Attachment{
						{
							Filename: "test.txt",
							MIMEType: "text/plain",
							B64Data:  "test",
						},
					},
				},
				Recipients: []string{"recipient@domain.com"},
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
			// Create a mock SparkPost client
			mockClient := &MockSparkPostClient{
				Response: tt.mockResponse,
				Error:    tt.mockError,
			}

			// Create a new instance of our service
			provider := &SparkpostProvider{
				Logger: logger,
				Client: mockClient,
			}

			err := provider.SendMail(tt.mail)

			// Assertions
			expectedContent, _ := tt.expectedSparkPostData.Content.(sparkpost.Content)
			actualContent, _ := mockClient.CalledWith.Content.(sparkpost.Content)
			assert.Equal(t, expectedContent.Text, actualContent.Text, "Unexpected Text content")
			assert.Equal(t, expectedContent.HTML, actualContent.HTML, "Unexpected HTML content")
			assert.Equal(t, expectedContent.From, actualContent.From, "Unexpected Sender")
			assert.Equal(t, expectedContent.Subject, actualContent.Subject, "Unexpected Subject")
			assert.Equal(t, expectedContent.Attachments, actualContent.Attachments, "Unexpected Attachments")
			assert.Equal(t, tt.expectedSparkPostData.Recipients, mockClient.CalledWith.Recipients, "Unexpected Recipients")
			assert.Equal(t, tt.expectedError, err, "Unexpected error")
		})
	}
}
