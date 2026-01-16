package email

import (
	"context"
	"email-sender/internal/config"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"time"
)

type Sender struct {
	cfg *config.Config
}

func NewSender(cfg *config.Config) *Sender {
	return &Sender{cfg: cfg}
}

func (s *Sender) SendWelcome(ctx context.Context, payload []byte) error {
	var event UserRegisteredEvent

	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("deserialize email payload: %w", err)
	}

	msg := buildWelcomeEmail(event)

	return s.sendSMTP(ctx, msg)
}

func (s *Sender) SendReport(ctx context.Context, payload []byte) error {
	var event DailyReportEvent

	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("deserialize email payload: %w", err)
	}

	msg := buildReportEmail(event)

	return s.sendSMTP(ctx, msg)
}

func (s *Sender) sendSMTP(ctx context.Context, msg Message) error {
	if msg.To == "" || msg.Subject == "" {
		return fmt.Errorf("invalid email payload")
	}

	smtpCfg := s.cfg.SmtpCfg
	addr := net.JoinHostPort(smtpCfg.Host, smtpCfg.Port)

	dialer := net.Dialer{
		Timeout: time.Duration(smtpCfg.TimeoutSec) * time.Second,
	}

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("smtp dial error: %w", err)
	}

	defer func() {
		if closeErr := conn.Close(); err != nil {
			log.Printf("conn close error: %v", closeErr)
		}
	}()

	client, err := smtp.NewClient(conn, smtpCfg.Host)
	if err != nil {
		return fmt.Errorf("smtp client error: %w", err)
	}

	defer func() {
		if quitErr := client.Quit(); quitErr != nil {
			log.Printf("client quit error: %v", quitErr)
		}
	}()

	if smtpCfg.Username != "" {
		auth := smtp.PlainAuth(
			"",
			smtpCfg.Username,
			smtpCfg.Password,
			smtpCfg.Host,
		)

		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth error: %w", err)
		}
	}

	if err = client.Mail(smtpCfg.From); err != nil {
		return fmt.Errorf("smtp mail from error: %w", err)
	}

	if err = client.Rcpt(msg.To); err != nil {
		return fmt.Errorf("smtp rcpt error: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data error: %w", err)
	}

	email := buildEmail(
		smtpCfg.From,
		msg.To,
		msg.Subject,
	)

	if _, err := w.Write([]byte(email)); err != nil {
		return fmt.Errorf("smtp write error: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close error: %w", err)
	}

	return nil
}

func buildWelcomeEmail(e UserRegisteredEvent) Message {
	return Message{
		To:      e.Email,
		Subject: fmt.Sprintf("Добро пожаловать, %s!", e.Email),
	}
}

func buildReportEmail(e DailyReportEvent) Message {
	return Message{
		To:      e.Email,
		Subject: fmt.Sprintf("Здравствуйте, %s! Сообщаем вам: %s", e.Email, e.Message),
	}
}

func buildEmail(from, to, subject string) string {
	return fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/plain; charset=\"utf-8\"\r\n"+
			"\r\n",
		from,
		to,
		subject,
	)
}
