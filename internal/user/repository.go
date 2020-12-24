package user

import (
	"context"
	"errors"

	apperr "go-api-template/internal/error"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
)

// Repository is data provider
type Repository interface {
	Create(ctx context.Context, u *User) (*User, error)
	Update(ctx context.Context, u *User) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id int) (*User, error)
}

var (
	errRepoUserAlreadyExists = apperr.New("repo", "user already exists", nil)
	errRepoUserNotFound      = apperr.New("repo", "user not found", nil)
	errRepoUnknown           = apperr.New("repo", "unkown error", nil)
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
			return nil, errRepoUserAlreadyExists.CloneWithInner(err)
		}

		return nil, errRepoUnknown.CloneWithInner(err)
	}

	return u, nil
}

func (r repo) Update(ctx context.Context, u *User) (*User, error) {
	_, err := r.db.ModelContext(ctx, u).WherePK().Update()
	if err != nil {
		r.logger.Debug().Err(err).Msg("")
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errRepoUserNotFound.CloneWithInner(err)
		}

		return nil, errRepoUnknown.CloneWithInner(err)
	}

	return u, nil
}

func (r repo) FindByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}

	err := r.db.ModelContext(ctx, u).Where("email = ?", email).First()
	if err != nil {
		r.logger.Debug().Err(err).Msg("")
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errRepoUserNotFound.CloneWithInner(err)
		}
		return nil, errRepoUnknown.CloneWithInner(err)
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
		return nil, errRepoUnknown.CloneWithInner(err)
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
