package mailer

import (
	"context"

	"github.com/rs/zerolog"
)

type fakeService struct {
	from     string
	fromName string
	logger   zerolog.Logger
}

func (s fakeService) Mail(ctx context.Context, email *Email) error {
	s.logger.Debug().Msgf("Mail sent! to: %v, subject: %s, body: %s", email.To, email.Subject, email.Body)
	return nil
}

// NewService returns a new fake service
func NewService(logger zerolog.Logger, from, formName string) Service {
	return &fakeService{
		logger:   logger,
		from:     from,
		fromName: formName,
	}
}
