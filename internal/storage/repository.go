package storage

import (
	"context"
	"errors"
	apperr "go-api-template/internal/error"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
)

var (
	errRepoUnknown           = apperr.New("repo", "unkown error", nil)
	errRepoFileNotFound      = apperr.New("repo", "file not found", nil)
	errRepoFileAlreadyExists = apperr.New("repo", "file already exists", nil)
)

// Repository is a storage service Repository
type Repository interface {
	Create(ctx context.Context, f *File) (*File, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	ExistsByIDAndUserID(ctx context.Context, id string, userID int) (bool, error)
	GetByID(ctx context.Context, id string) (*File, error)
	GetByIDAndUserID(ctx context.Context, id string, userID int) (*File, error)
	Delete(ctx context.Context, id string) error
	SetUploadStatusByID(ctx context.Context, id string, uploadStatus UploadStatus) (*File, error)
	SetUploadStatusByIDAndUserID(ctx context.Context, id string, userID int, uploadStatus UploadStatus) (*File, error)
	SetUploadStatusAndPublicURLByID(ctx context.Context, id string, uploadStatus UploadStatus, publicURL string) (*File, error)
	SetUploadStatusAndPublicURLByIDAndUserID(ctx context.Context, id string, userID int, uploadStatus UploadStatus, publicURL string) (*File, error)
}

type repo struct {
	logger zerolog.Logger
	db     *pg.DB
}

func (r repo) Create(ctx context.Context, f *File) (*File, error) {
	_, err := r.db.ModelContext(ctx, f).Insert()
	if err != nil {
		return nil, errRepoUnknown.CloneWithInner(err)
	}

	return f, nil
}

func (r repo) ExistsByID(ctx context.Context, id string) (bool, error) {
	exists, err := r.db.ModelContext(ctx, &File{ID: id}).WherePK().Exists()
	if err != nil {
		return false, errRepoUnknown.CloneWithInner(err)
	}

	return exists, nil
}

func (r repo) ExistsByIDAndUserID(ctx context.Context, id string, userID int) (bool, error) {
	exists, err := r.db.
		ModelContext(ctx, &File{ID: id}).
		Where("id = ? and user_id = ?", id, userID).
		Exists()
	if err != nil {
		return false, errRepoUnknown.CloneWithInner(err)
	}

	return exists, nil
}

func (r repo) GetByID(ctx context.Context, id string) (*File, error) {
	f := &File{}
	err := r.db.ModelContext(ctx, f).Where("id = ?", id).First()
	if err != nil {
		r.logger.Debug().Err(err).Msg("")
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errRepoFileNotFound.CloneWithInner(err)
		}
		return nil, errRepoUnknown.CloneWithInner(err)
	}

	return f, nil
}

func (r repo) GetByIDAndUserID(ctx context.Context, id string, userID int) (*File, error) {
	f := &File{}
	err := r.db.
		ModelContext(ctx, f).
		Where("id = ? and user_id = ?", id, userID).
		First()
	if err != nil {
		r.logger.Debug().Err(err).Msg("")
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errRepoFileNotFound.CloneWithInner(err)
		}
		return nil, errRepoUnknown.CloneWithInner(err)
	}

	return f, nil
}

func (r repo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ModelContext(ctx, &File{ID: id}).WherePK().Delete()
	if err != nil {
		r.logger.Debug().Err(err).Msg("")
		if errors.Is(err, pg.ErrNoRows) {
			return errRepoFileNotFound.CloneWithInner(err)
		}
		return errRepoUnknown.CloneWithInner(err)
	}

	return nil
}

func (r repo) SetUploadStatusByID(ctx context.Context, id string, uploadStatus UploadStatus) (*File, error) {
	f := &File{}
	_, err := r.db.
		ModelContext(ctx, f).
		Set("upload_status = ?", uploadStatus).
		Where("id = ?", id).
		Returning("*").
		Update()
	if err != nil {
		r.logger.Debug().Err(err).Msg("")
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errRepoFileNotFound.CloneWithInner(err)
		}
		return nil, errRepoUnknown.CloneWithInner(err)
	}

	return f, nil
}

func (r repo) SetUploadStatusByIDAndUserID(ctx context.Context, id string, userID int, uploadStatus UploadStatus) (*File, error) {
	f := &File{}
	_, err := r.db.
		ModelContext(ctx, f).
		Set("upload_status = ?", uploadStatus).
		Where("id = ? and user_id = ?", id, userID).
		Returning("*").
		Update()
	if err != nil {
		r.logger.Debug().Err(err).Msg("")
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errRepoFileNotFound.CloneWithInner(err)
		}
		return nil, errRepoUnknown.CloneWithInner(err)
	}

	return f, nil
}

func (r repo) SetUploadStatusAndPublicURLByID(ctx context.Context, id string, uploadStatus UploadStatus, publicURL string) (*File, error) {
	f := &File{}
	_, err := r.db.
		ModelContext(ctx, f).
		Set("upload_status = ? and public_url", uploadStatus, publicURL).
		Where("id = ?", id).
		Returning("*").
		Update()
	if err != nil {
		r.logger.Debug().Err(err).Msg("")
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errRepoFileNotFound.CloneWithInner(err)
		}
		return nil, errRepoUnknown.CloneWithInner(err)
	}

	return f, nil
}

func (r repo) SetUploadStatusAndPublicURLByIDAndUserID(ctx context.Context, id string, userID int, uploadStatus UploadStatus, publicURL string) (*File, error) {
	f := &File{}
	_, err := r.db.
		ModelContext(ctx, f).
		Set("upload_status = ? and public_url = ?", uploadStatus, publicURL).
		Where("id = ? and user_id = ?", id, userID).
		Returning("*").
		Update()
	if err != nil {
		r.logger.Debug().Err(err).Msg("")
		if errors.Is(err, pg.ErrNoRows) {
			return nil, errRepoFileNotFound.CloneWithInner(err)
		}
		return nil, errRepoUnknown.CloneWithInner(err)
	}

	return f, nil
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
