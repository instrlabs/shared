package email

import (
	"fmt"
	"net/smtp"

	initx "github.com/instrlabs/shared/init"
)

var sendMail = smtp.SendMail

func SendEmail(emailTo, subject, body string) error {
	from := initx.GetEnv("EMAIL_FROM", "")
	host := initx.GetEnv("SMTP_HOST", "")
	port := initx.GetEnv("SMTP_PORT", "")
	username := initx.GetEnv("SMTP_USERNAME", "")
	password := initx.GetEnv("SMTP_PASSWORD", "")

	if from == "" || host == "" || port == "" || username == "" || password == "" {
		return fmt.Errorf("missing SMTP configuration in environment")
	}

	// Build minimal MIME message headers
	headers := map[string]string{
		"From":         from,
		"To":           emailTo,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/plain; charset=\"UTF-8\"",
	}
	var msg string
	for k, v := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + body

	addr := host + ":" + port
	auth := smtp.PlainAuth("", username, password, host)

	if err := sendMail(addr, auth, from, []string{emailTo}, []byte(msg)); err != nil {
		return fmt.Errorf("sendmail failed: %w", err)
	}
	return nil
}
