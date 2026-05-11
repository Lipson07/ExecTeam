// internal/service/email.go
package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"backend/config"
)

type EmailService struct {
	from     string
	password string
}

func NewEmailService(cfg *config.Config) *EmailService {
	fmt.Printf("EMAIL INIT: from=%s pass=%s\n", cfg.SMTPUsername, cfg.SMTPPassword)
	return &EmailService{
		from:     cfg.SMTPUsername,
		password: cfg.SMTPPassword,
	}
}

func (s *EmailService) SendVerificationCode(to string, code string) error {
	fmt.Printf("SEND START: to=%s code=%s from=%s pass=%s\n", to, code, s.from, s.password)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Код\r\n\r\nКод: %s\r\n", s.from, to, code)

	tlsConfig := &tls.Config{ServerName: "smtp.yandex.ru"}

	conn, err := tls.Dial("tcp", "smtp.yandex.ru:465", tlsConfig)
	fmt.Printf("DIAL: err=%v\n", err)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, "smtp.yandex.ru")
	fmt.Printf("CLIENT: err=%v\n", err)
	if err != nil {
		return err
	}
	defer client.Quit()

	auth := smtp.PlainAuth("", s.from, s.password, "smtp.yandex.ru")
	err = client.Auth(auth)
	fmt.Printf("AUTH: err=%v\n", err)
	if err != nil {
		return err
	}

	err = client.Mail(s.from)
	fmt.Printf("MAIL: err=%v\n", err)
	if err != nil {
		return err
	}

	err = client.Rcpt(to)
	fmt.Printf("RCPT: err=%v\n", err)
	if err != nil {
		return err
	}

	w, err := client.Data()
	fmt.Printf("DATA: err=%v\n", err)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(msg))
	fmt.Printf("WRITE: err=%v\n", err)

	err = w.Close()
	fmt.Printf("CLOSE: err=%v\n", err)

	fmt.Printf("EMAIL SENT: to=%s code=%s\n", to, code)
	return nil
}
