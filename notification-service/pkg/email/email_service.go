package email

import (
	"context"
	"fmt"
	"html/template"
	"micro-warehouse/notification-service/configs"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/gomail.v2"
)

type EmailServiceInterface interface {
	SendWelcomeEmail(ctx context.Context, payload EmailPayload) error
	SendCustomEmail(ctx context.Context, to, subject, body string) error
}

type EmailPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
	UserId   string `json:"user_id"`
	Name     string `json:"name"`
}

type emailService struct {
	cfg configs.Config
}

// SendCustomEmail implements [EmailServiceInterface].
func (e *emailService) SendCustomEmail(ctx context.Context, to string, subject string, body string) error {
	if e.cfg.Email.Host == "" || e.cfg.Email.User == "" || e.cfg.Email.Password == "" {
		log.Errorf("[EmailService] SendCustomEmail - 1: %v", "email configuration is incomplete")
		return fmt.Errorf("email configuration is incomplete: Host=%s, User=%s", e.cfg.Email.Host, e.cfg.Email.User)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", e.cfg.Email.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(e.cfg.Email.Host, e.cfg.Email.Port, e.cfg.Email.User, e.cfg.Email.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Errorf("[EmailService] SendCustomEmail - 2: %v", err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// SendWelcomeEmail implements [EmailServiceInterface].
func (e *emailService) SendWelcomeEmail(ctx context.Context, payload EmailPayload) error {
	subject := "Selamat Datang di Warehouse Management"
	htmlTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Selamat Datang</title>
		<style>
			body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
			.container { max-width: 600px; margin: 0 auto; padding: 20px; }
			.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
			.content { padding: 20px; background-color: #f9f9f9; }
			.footer { text-align: center; padding: 20px; background-color: #f4f4f4; color: #666; font-size: 12px; }
			.button { display: inline-block; padding: 10px 20px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h1>Selamat Datang</h1>
			</div>
			<div class="content">
				<p>Halo {{.Name}},</p>
				<p>Selamat datang di Warehouse Management. Berikut adalah detail akun Anda:</p>
				<p><strong>Email:</strong> {{.Email}}</p>
				<p><strong>Password:</strong> {{.Password}}</p>
				<p>Anda dapat login ke aplikasi dengan menggunakan email dan password tersebut. Silakan untuk mengganti password setelah pertama kali login.</p>
				<a href="{{.AppURL}}" class="button">Login ke Aplikasi</a>
			</div>
			<div class="footer">
				<p>Email ini dikirim secara otomatis, mohon jangan membalas email ini.</p>
				<p>&copy; 2026 Warehouse Management. All rights reserved.</p>
			</div>
		</div>
	</body>
	</html>
	`

	tmpl, err := template.New("welcome").Parse(htmlTemplate)
	if err != nil {
		log.Errorf("[EmailService] SendWelcomeEmail - 1: %v", err)
		return fmt.Errorf("failed to parse email template: %v", err)
	}

	var body strings.Builder
	err = tmpl.Execute(&body, payload)
	if err != nil {
		log.Errorf("[EmailService] SendWelcomeEmail - 2: %v", err)
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	err = e.SendCustomEmail(ctx, payload.Email, subject, body.String())
	if err != nil {
		log.Errorf("[EmailService] SendWelcomeEmail - 3: %v", err)
		return fmt.Errorf("failed to send welcome email: %v", err)
	}

	return nil
}

func NewEmailService(cfg configs.Config) EmailServiceInterface {
	return &emailService{
		cfg: cfg,
	}
}
