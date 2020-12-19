package user

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
)

// Repository is data provider
type Repository interface {
	Create(ctx context.Context, u *User) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id int) (*User, error)
}

var (
	errRepoUserAlreadyExists = errors.New("user already exists")
	errRepoUserNotFound      = errors.New("user not found")
)

type repo struct {
	logger zerolog.Logger
	db     *pg.DB
}

func (r repo) Create(ctx context.Context, u *User) (*User, error) {
	_, err := r.db.ModelContext(ctx, u).Insert()
	if err != nil {
		r.logger.Debug().Err(err).Msg("")
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			return nil, errRepoUserAlreadyExists
		}

		return nil, err
	}

	// r.logger.Debug().Int("user_id", u.ID).Str("email", u.Email).Msg("user created")

	return u, nil
}

func (r repo) FindByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}

	err := r.db.ModelContext(ctx, u).Where("email = ?", email).First()
	if err != nil {
		r.logger.Debug().Err(err).Msg("")
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errRepoUserNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r repo) FindByID(ctx context.Context, id int) (*User, error) {
	u := &User{}

	err := r.db.ModelContext(ctx, u).Where("id = ?", id).First()
	if err != nil {
		r.logger.Debug().Err(err).Msg("")
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errRepoUserNotFound
		}
		return nil, err
	}

	return u, nil
}

// NewRepository creates a new repository
func NewRepository(
	logger zerolog.Logger,
	db *pg.DB,
) Repository {
	return &repo{
		logger: logger,
		db:     db,
	}
}
