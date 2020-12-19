package user

import (
	"context"
	"errors"

	pass "github.com/dev681999/go-pass"
	"github.com/rs/zerolog"
)

// Service is a service provider
type Service interface {
	Create(ctx context.Context, u *User) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id int) (*User, error)
}

// Errors that can occur in the service
var (
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInternalService   = errors.New("internal service error")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type service struct {
	logger zerolog.Logger
	repo   Repository
	pass.Hash
}

func (s service) Create(ctx context.Context, u *User) (*User, error) {
	password, err := s.Generate(u.Password)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		err = ErrInvalidPassword
		return nil, err
	}

	u.Password = password

	u, err = s.repo.Create(ctx, u)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if errors.Is(err, errRepoUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
		}

		return nil, ErrInternalService
	}

	return u, nil
}

func (s service) FindByEmail(ctx context.Context, email string) (*User, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if errors.Is(err, errRepoUserNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, ErrInternalService
	}

	u.Password = ""

	return u, nil
}

func (s service) FindByID(ctx context.Context, id int) (*User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if errors.Is(err, errRepoUserNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, ErrInternalService
	}

	u.Password = ""

	return u, nil
}

// NewService creates a new service
func NewService(
	logger zerolog.Logger,
	repo Repository,
) Repository {
	return &service{
		logger: logger,
		repo:   repo,
	}
}
