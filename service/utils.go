package service

import (
	"fmt"
	"github.com/gugabfigueiredo/dream-mail-go/models"
	"path/filepath"
	"strings"
)

func buildSMTPMessage(mail *models.Mail, tos []string) []byte {
	//build message from mail
	msg := fmt.Sprintf("From: %s\r\n", mail.From.Addr)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(tos, ","))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)

	msg += "MIME-Version: 1.0\r\n"
	msg += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"dmailboundary\"\r\n")

	if mail.Text != "" {
		msg += "\r\n--dmailboundary\r\n"
		msg += "Content-Type: text/plain; charset=\"utf-8\"\r\n"
		msg += fmt.Sprintf("\r\n%s\r\n", mail.Text)
	}

	if mail.HTML != "" {
		msg += "\r\n--dmailboundary\r\n"
		msg += "Content-Type: text/html; charset=\"utf-8\"\r\n"
		msg += fmt.Sprintf("\r\n%s\r\n", mail.HTML)
	}

	for _, attachment := range mail.Attachments {
		_, fileName := filepath.Split(attachment.Name)

		msg += "\r\n--dmailboundary\r\n"
		msg += fmt.Sprintf("Content-Type: %s; charset=\"utf-8\"\r\n", attachment.Type)
		msg += fmt.Sprintf("Content-Disposition: attachment;filename=\"%s\"\r\n", fileName)
		msg += fmt.Sprintf("\r\n%s\r\n", attachment.Data)
	}

	msg += "\r\n--dmailboundary\r\n"

	return []byte(msg)
}
