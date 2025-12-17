package repository

import (
	"context"
	"fmt"
	"net/smtp"

	"good-todo-go/internal/domain/repository"
	"good-todo-go/internal/infrastructure/environment"
)

type AuthRepository struct {
	env *environment.Environment
}

func NewAuthRepository(env *environment.Environment) repository.IAuthRepository {
	return &AuthRepository{env: env}
}

func (r *AuthRepository) SendVerificationEmail(ctx context.Context, email, token string) error {
	from := "noreply@good-todo-go.local"
	to := []string{email}
	subject := "Please verify your email"
	body := fmt.Sprintf(`Hello,

Please verify your email by clicking the link below:

http://localhost:3000/verify-email?token=%s

This link will expire in 24 hours.

Best regards,
Good Todo Go Team`, token)

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, email, subject, body))

	addr := fmt.Sprintf("%s:%s", r.env.SMTPHost, r.env.SMTPPort)
	err := smtp.SendMail(addr, nil, from, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}
