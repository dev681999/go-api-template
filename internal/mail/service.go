package mail

import (
	"context"
	"fmt"
	"go-api-template/internal/mailer"

	"github.com/rs/zerolog"
)

// Service is a mail service provider
type Service interface {
	SendWelcomeMail(ctx context.Context, to string, activationURL string) error
	SendPasswordResetLink(ctx context.Context, to string, resetLink string) error
}

type service struct {
	logger zerolog.Logger
	ms     mailer.Service
}

func (s service) SendWelcomeMail(ctx context.Context, to string, activationURL string) error {
	subject := "Welcome"
	body := fmt.Sprintf("Welcome! Please click here to activate your account. %s", activationURL)

	return s.ms.Mail(ctx, &mailer.Email{
		To:      []string{to},
		Subject: subject,
		Body:    []byte(body),
	})
}

func (s service) SendPasswordResetLink(ctx context.Context, to string, resetLink string) error {
	subject := "Password reset"
	body := fmt.Sprintf("Please click below link to reset your password <br/> %s", resetLink)

	return s.ms.Mail(ctx, &mailer.Email{
		To:      []string{to},
		Subject: subject,
		Body:    []byte(body),
	})
}

// NewService returns a new mail service
func NewService(logger zerolog.Logger, ms mailer.Service) Service {
	return &service{
		logger: logger,
		ms:     ms,
	}
}
