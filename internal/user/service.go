package user

import (
	"context"
	"fmt"
	apperr "go-api-template/internal/error"
	"go-api-template/internal/mail"
	"net/url"
	"time"

	pass "github.com/dev681999/go-pass"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog"
)

// Service is a service provider
type Service interface {
	Create(ctx context.Context, u *User) (*User, *apperr.Error)
	Update(ctx context.Context, userID int, update *Update) (*User, *apperr.Error)

	FindByID(ctx context.Context, id int) (*User, *apperr.Error)
	FindByEmail(ctx context.Context, email string) (*User, *apperr.Error)

	ValidateByEmailAndPassword(ctx context.Context, email, password string) (*User, *apperr.Error)

	GenerateToken(ctx context.Context, u *User, claimType ClaimType) (string, *apperr.Error)
	ParseTokenWithClaims(ctx context.Context, auth string) (interface{}, bool, error)
}

// Errors that can occur in the service
var (
	ErrInvalidPassword   = apperr.New("service", "invalid password", 0, nil)
	ErrInternalService   = apperr.New("service", "internal service error", 1, nil)
	ErrUserAlreadyExists = apperr.New("service", "user already exists", 2, nil)
	ErrUserNotFound      = apperr.New("service", "user not found", 3, nil)
	ErrWrongCredentials  = apperr.New("service", "wrong credentials", 4, nil)
	ErrUserNotActive     = apperr.New("service", "user not active", 5, nil)
	ErrInvalidToken      = apperr.New("service", "invalid token", 6, nil)
)

type service struct {
	logger zerolog.Logger
	repo   Repository
	ms     mail.Service

	jwtKey        string
	activationURL string

	pass.Hash
}

func (s service) Create(ctx context.Context, u *User) (*User, *apperr.Error) {
	{
		password, err := s.Generate(u.Password)
		if err != nil {
			s.logger.Debug().Err(err).Msg("")
			return nil, ErrInvalidPassword.CloneWithInner(err)
		}

		u.Password = password
	}

	u, err := s.repo.Create(ctx, u)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if err.IsMatchesCode(errRepoUserAlreadyExists) {
			return nil, ErrUserAlreadyExists.CloneWithInner(err)
		}

		return nil, ErrInternalService.CloneWithInner(err)
	}

	token, err := s.GenerateToken(ctx, u, ClaimTypeActivation)
	if err != nil {
		return nil, ErrInternalService.CloneWithInner(err)
	}

	activationURL := ""
	{
		pURL, err := url.Parse(s.activationURL)
		if err != nil {
			return nil, ErrInternalService.CloneWithInner(err)
		}
		q := pURL.Query()
		q.Set("token", token)
		pURL.RawQuery = q.Encode()

		activationURL = pURL.String()
	}

	if err := s.ms.SendWelcomeMail(ctx, u.Email, activationURL); err != nil {
		s.logger.Debug().Err(err).Msg("")
		return nil, ErrInternalService.CloneWithInner(err)
	}

	return u, nil
}

func (s service) Update(ctx context.Context, userID int, update *Update) (*User, *apperr.Error) {
	u, err := s.findByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	u = populateUpdateInUser(u, update)

	u, err = s.repo.Update(ctx, u)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if err.IsMatchesCode(errRepoUserNotFound) {
			return nil, ErrUserNotFound.CloneWithInner(err)
		}

		return nil, ErrInternalService.CloneWithInner(err)
	}

	return u, nil
}

func (s service) FindByEmail(ctx context.Context, email string) (*User, *apperr.Error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if err.IsMatchesCode(errRepoUserNotFound) {
			return nil, ErrUserNotFound.CloneWithInner(err)
		}

		return nil, ErrInternalService.CloneWithInner(err)
	}

	u.Password = ""

	return u, nil
}

func (s service) ValidateByEmailAndPassword(ctx context.Context, email, password string) (*User, *apperr.Error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if err.IsMatchesCode(errRepoUserNotFound) {
			return nil, ErrUserNotFound.CloneWithInner(err)
		}

		return nil, ErrInternalService.CloneWithInner(err)
	}

	if err := s.Compare(u.Password, password); err != nil {
		return nil, ErrWrongCredentials.CloneWithInner(err)
	}

	u.Password = ""

	return u, nil
}

func (s service) findByID(ctx context.Context, id int) (*User, *apperr.Error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if err.IsMatchesCode(errRepoUserNotFound) {
			return nil, ErrUserNotFound.CloneWithInner(err)
		}

		return nil, ErrInternalService.CloneWithInner(err)
	}

	return u, nil
}

func (s service) FindByID(ctx context.Context, id int) (*User, *apperr.Error) {
	u, err := s.findByID(ctx, id)
	if err != nil {
		return nil, err
	}

	u.Password = ""

	return u, nil
}

func (s service) GenerateToken(ctx context.Context, u *User, claimType ClaimType) (string, *apperr.Error) {
	var expiresAt int64

	switch claimType {
	case ClaimTypeActivation:
		expiresAt = time.Now().Add(time.Hour * 24).Unix()
	default:
		expiresAt = time.Now().Add(time.Hour * 24 * 365).Unix()
	}

	claims := &Claims{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		UserID:    u.ID,
		Type:      claimType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(s.jwtKey))
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		return "", ErrInternalService.CloneWithInner(err)
	}

	return t, nil
}

// Defaults
const (
	AlgorithmHS256 = "HS256"
)

func (s service) ParseTokenWithClaims(ctx context.Context, auth string) (interface{}, bool, error) {
	claims := &Claims{}

	{
		token, err := jwt.ParseWithClaims(auth, claims, func(t *jwt.Token) (interface{}, error) {
			// Check the signing method
			if t.Method.Alg() != AlgorithmHS256 {
				return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
			}

			return []byte(s.jwtKey), nil
		})
		if err != nil {
			return nil, false, err
		}

		if !token.Valid {
			return nil, false, ErrInvalidToken
		}

		if err := claims.Valid(); err != nil {
			return nil, false, ErrInvalidToken.CloneWithInner(err)
		}
	}

	u, err := s.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, false, err
	}

	if !u.Active {
		return nil, false, ErrUserNotActive
	}

	return claims, true, nil
}

// NewService creates a new service
func NewService(
	logger zerolog.Logger,
	repo Repository,
	ms mail.Service,
	jwtKey string,
	activationURL string,
) Service {
	return &service{
		logger:        logger,
		repo:          repo,
		ms:            ms,
		jwtKey:        jwtKey,
		activationURL: activationURL,
	}
}
