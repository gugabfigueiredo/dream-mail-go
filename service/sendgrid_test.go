package service

import (
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	sgMail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"testing"
)

type MockSendClient struct {
	status int
	err    error
}

func (c *MockSendClient) Send(m *sgMail.SGMailV3) (*rest.Response, error) {
	if c.err != nil {
		return nil, c.err
	}

	// Perform any necessary assertions on the SGMailV3 object before returning a mock response
	// ...

	// Return a mock response
	return &rest.Response{
		StatusCode: c.status,
	}, nil
}

func TestSendEmail(t *testing.T) {
	// Set up your SendGrid API key and client
	apiKey := "YOUR_API_KEY"
	client := sendgrid.NewSendClient(apiKey)

	// Define the test cases
	testCases := []struct {
		name               string
		fromName           string
		fromEmail          string
		toName             string
		toEmail            string
		subject            string
		contentType        string
		content            string
		attachments        []*sgMail.Attachment
		personalizations   []*sgMail.Personalization
		expectedStatusCode int
	}{
		{
			name: "e-mail with attachments and personalizations - should call Send and return status ok",
		},
		{
			name: "e-mail with missing key params - should not call Send and return error",
		},
		{
			name: "e-mail",
		},
		{
			name:        "e-mail with personalizations - should call Send and return status ok",
			fromName:    "Sender Name",
			fromEmail:   "sender@example.com",
			toName:      "Recipient Name",
			toEmail:     "recipient@example.com",
			subject:     "Test Email",
			contentType: "text/plain",
			content:     "Hello, this is a test email!",
			attachments: []*sgMail.Attachment{
				{
					Content:     "Attachment 1 content",
					Type:        "application/pdf",
					Filename:    "attachment1.pdf",
					Disposition: "attachment",
					ContentID:   "attachment1",
				},
				{
					Content:     "Attachment 2 content",
					Type:        "image/jpeg",
					Filename:    "attachment2.jpg",
					Disposition: "attachment",
					ContentID:   "attachment2",
				},
			},
			personalizations: []*sgMail.Personalization{
				{
					To: []*sgMail.Email{
						sgMail.NewEmail("Recipient Name", "recipient@example.com"),
					},
					Subject: "Personalized Subject",
					Headers: map[string]string{
						"X-Custom-Header": "Custom Value",
					},
					Substitutions: map[string]string{
						"%name%": "John Doe",
						"%city%": "New York",
					},
				},
			},
			expectedStatusCode: 202,
		},
		{
			name:               "Email without Attachments or Personalizations",
			fromName:           "Sender Name",
			fromEmail:          "sender@example.com",
			toName:             "Recipient Name",
			toEmail:            "recipient@example.com",
			subject:            "Test Email",
			contentType:        "text/plain",
			content:            "Hello, this is a test email!",
			attachments:        []*sgMail.Attachment{},
			personalizations:   []*sgMail.Personalization{},
			expectedStatusCode: 202,
		},
		// Add more test cases here...
	}

	// Iterate over the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the email content
			from := sgMail.NewEmail(tc.fromName, tc.fromEmail)
			to := sgMail.NewEmail(tc.toName, tc.toEmail)
			content := sgMail.NewContent(tc.contentType, tc.content)

			// Create the email message
			message := sgMail.NewV3MailInit(from, tc.subject, to, content)

			// Add attachments
			for _, attachment := range tc.attachments {
				message.AddAttachment(attachment)
			}

			// Add personalizations
			for _, personalization := range tc.personalizations {
				message.AddPersonalizations(personalization)
			}

			// Send the email
			response, err := client.Send(message)
			if err != nil {
				t.Errorf("Error sending email: %v", err)
			}

			// Check the response status code
			if response.StatusCode != 202 {
				t.Errorf("Unexpected status code: %d", response.StatusCode)
			}
		})
	}
}
