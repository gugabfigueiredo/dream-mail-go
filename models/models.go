package models

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type Email struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type Attachment struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Data string `json:"data"`
}

type Mail struct {
	ID          string       `json:"id"`
	From        Email        `json:"from"`
	To          []Email      `json:"to"`
	Subject     string       `json:"subject"`
	Text        string       `json:"text"`
	HTML        string       `json:"html"`
	Attachments []Attachment `json:"attachments"`
}

func NewMail(from Email, subject, text string) *Mail {
	return &Mail{
		From:    from,
		Subject: subject,
		Text:    text,
	}
}

func (m *Mail) AddTo(to ...Email) {
	for _, email := range to {
		m.To = append(m.To, email)
	}
}

func (m *Mail) AddAttachment(name, mimeType string, data string) {
	m.Attachments = append(m.Attachments, Attachment{
		Name: name,
		Type: mimeType,
		Data: data,
	})
}

func (m *Mail) Validate() (bool, error) {

	// validate email
	if m.From.Addr == "" {
		return false, errors.New("missing sender")
	}

	if len(m.To) == 0 {
		return false, errors.New("missing recipient")
	}

	for _, recipient := range m.To {
		if recipient.Addr == "" {
			return false, errors.New("missing recipient address")
		}
	}

	if m.Subject == "" {
		return false, errors.New("missing subject")
	}

	return true, nil
}

type Provider struct {
	ID     int             `json:"id"`
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

type App struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Providers []int  `json:"providers"`
	APIKey    string `json:"api_key"`
}
