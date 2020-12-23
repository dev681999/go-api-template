package mailer

import "context"

// Service is a mailer service provider
type Service interface {
	Mail(ctx context.Context, email *Email) error
}

// Email is a email that needs to be sent
type Email struct {
	To      []string
	Subject string
	Body    []byte
}
