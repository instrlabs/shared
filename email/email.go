package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	initx "github.com/instrlabs/shared/init"
)

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

	var client *smtp.Client
	if port == "465" {
		// Implicit TLS
		conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host})
		if err != nil {
			return fmt.Errorf("tls dial failed: %w", err)
		}
		defer conn.Close()

		c, err := smtp.NewClient(conn, host)
		if err != nil {
			return fmt.Errorf("smtp new client failed: %w", err)
		}
		client = c
	} else {
		// Plain connect then STARTTLS if available
		c, err := smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("smtp dial failed: %w", err)
		}
		client = c
		// Upgrade to TLS if server supports it
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(&tls.Config{ServerName: host}); err != nil {
				_ = client.Close()
				return fmt.Errorf("starttls failed: %w", err)
			}
		}
	}
	defer client.Quit()

	// Authenticate if supported
	if ok, _ := client.Extension("AUTH"); ok {
		if err := client.Auth(auth); err != nil {
			_ = client.Close()
			return fmt.Errorf("smtp auth failed: %w", err)
		}
	}

	if err := client.Mail(from); err != nil {
		_ = client.Close()
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}
	if err := client.Rcpt(emailTo); err != nil {
		_ = client.Close()
		return fmt.Errorf("RCPT TO failed: %w", err)
	}
	wc, err := client.Data()
	if err != nil {
		_ = client.Close()
		return fmt.Errorf("DATA failed: %w", err)
	}
	if _, err := wc.Write([]byte(msg)); err != nil {
		_ = wc.Close()
		_ = client.Close()
		return fmt.Errorf("write message failed: %w", err)
	}
	_ = wc.Close()
	return nil
}
