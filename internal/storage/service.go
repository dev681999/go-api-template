package storage

import (
	"context"
	"errors"
	"fmt"
	apperr "go-api-template/internal/error"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Service is a storage service
type Service interface {
	Create(ctx context.Context, f *File) (*File, string, error)
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

// Errors that can occur in the service
var (
	ErrInternalService   = apperr.New("service", "internal service error", nil)
	ErrFileAlreadyExists = apperr.New("service", "file already exists", nil)
	ErrFileNotFound      = apperr.New("service", "file not found", nil)
	ErrFailPresignedURL  = apperr.New("service", "failed to create presigned url", nil)
)

type service struct {
	logger   zerolog.Logger
	repo     Repository
	s3Client *s3.S3
}

func (s service) Create(ctx context.Context, f *File) (*File, string, error) {
	f.ID = uuid.New().String()
	f.Bucket = fmt.Sprintf("user-%v", f.UserID)
	f.ObjectKey = fmt.Sprintf("%v-%v", f.ID, f.FileName)
	f.UploadStatus = UploadStatusPending

	req, _ := s.s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: &f.Bucket,
		Key:    &f.ObjectKey,
	})

	urlStr, err := req.Presign(5 * time.Minute)
	if err != nil {
		return nil, "", ErrFailPresignedURL.CloneWithInner(err)
	}

	f, err = s.repo.Create(ctx, f)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if errors.Is(err, errRepoFileAlreadyExists) {
			return nil, "", ErrFileAlreadyExists.CloneWithInner(err)
		}

		return nil, "", ErrInternalService.CloneWithInner(err)
	}

	return f, urlStr, nil
}

func (s service) ExistsByID(ctx context.Context, id string) (bool, error) {
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		return false, ErrInternalService.CloneWithInner(err)
	}

	return exists, nil
}

func (s service) ExistsByIDAndUserID(ctx context.Context, id string, userID int) (bool, error) {
	exists, err := s.repo.ExistsByIDAndUserID(ctx, id, userID)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		return false, ErrInternalService.CloneWithInner(err)
	}

	return exists, nil
}

func (s service) GetByID(ctx context.Context, id string) (*File, error) {
	f, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if errors.Is(err, errRepoFileNotFound) {
			return nil, ErrFileNotFound.CloneWithInner(err)
		}

		return nil, ErrInternalService.CloneWithInner(err)
	}

	return f, nil
}

func (s service) GetByIDAndUserID(ctx context.Context, id string, userID int) (*File, error) {
	f, err := s.repo.GetByIDAndUserID(ctx, id, userID)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if errors.Is(err, errRepoFileNotFound) {
			return nil, ErrFileNotFound.CloneWithInner(err)
		}

		return nil, ErrInternalService.CloneWithInner(err)
	}

	return f, nil
}

func (s service) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if errors.Is(err, errRepoFileNotFound) {
			return ErrFileNotFound.CloneWithInner(err)
		}

		return ErrInternalService.CloneWithInner(err)
	}

	return nil
}

func (s service) SetUploadStatusByID(ctx context.Context, id string, uploadStatus UploadStatus) (*File, error) {
	f, err := s.repo.SetUploadStatusByID(ctx, id, uploadStatus)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if errors.Is(err, errRepoFileNotFound) {
			return nil, ErrFileNotFound.CloneWithInner(err)
		}

		return nil, ErrInternalService.CloneWithInner(err)
	}

	return f, nil
}

func (s service) SetUploadStatusByIDAndUserID(ctx context.Context, id string, userID int, uploadStatus UploadStatus) (*File, error) {
	f, err := s.repo.SetUploadStatusByIDAndUserID(ctx, id, userID, uploadStatus)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if errors.Is(err, errRepoFileNotFound) {
			return nil, ErrFileNotFound.CloneWithInner(err)
		}

		return nil, ErrInternalService.CloneWithInner(err)
	}

	return f, nil
}

func (s service) SetUploadStatusAndPublicURLByID(ctx context.Context, id string, uploadStatus UploadStatus, publicURL string) (*File, error) {
	f, err := s.repo.SetUploadStatusAndPublicURLByID(ctx, id, uploadStatus, publicURL)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if errors.Is(err, errRepoFileNotFound) {
			return nil, ErrFileNotFound.CloneWithInner(err)
		}

		return nil, ErrInternalService.CloneWithInner(err)
	}

	return f, nil
}

func (s service) SetUploadStatusAndPublicURLByIDAndUserID(ctx context.Context, id string, userID int, uploadStatus UploadStatus, publicURL string) (*File, error) {
	f, err := s.repo.SetUploadStatusAndPublicURLByIDAndUserID(ctx, id, userID, uploadStatus, publicURL)
	if err != nil {
		s.logger.Debug().Err(err).Msg("")
		if errors.Is(err, errRepoFileNotFound) {
			return nil, ErrFileNotFound.CloneWithInner(err)
		}

		return nil, ErrInternalService.CloneWithInner(err)
	}

	return f, nil
}

// NewService returns a new service
func NewService(
	logger zerolog.Logger,
	repo Repository,
	s3Client *s3.S3,
) Service {
	return &service{
		logger:   logger,
		repo:     repo,
		s3Client: s3Client,
	}
}
