package service

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	from     string
	password string
}

func NewEmailService(from string, password string) *EmailService {
	return &EmailService{from: from, password: password}
}

func (s *EmailService) SendVerificationCode(email, code string) error {
	subject := "Код подтверждения"
	body := fmt.Sprintf(`Привет!

Твой код подтверждения: %s

Код действителен 10 минут. Если ты не регистрировался — просто проигнорируй это письмо.`, code)

	msg := []byte("From: " + s.from + "\r\n" +
		"To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body)

	auth := smtp.PlainAuth("", s.from, s.password, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, s.from, []string{email}, msg)
	if err != nil {
		return fmt.Errorf("send verification email: %w", err)
	}
	return nil
}
