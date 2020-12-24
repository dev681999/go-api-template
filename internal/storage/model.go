package storage

import (
	"context"
	"time"
)

// UploadStatus is status of the upload
type UploadStatus uint

// UploafStatuses
const (
	UploadStatusPending UploadStatus = iota
	UploadStatusUploaded
	UploadStatusFailed
)

// FileType is type of file
type FileType uint

// FileTypes available
const (
	FileTypeImage FileType = iota
	FileTypeVideo
)

// File is a file in the storage
type File struct {
	tableName struct{} `pg:"users,alias:users"`

	ID           string       `pg:",pk" json:"-"`
	UserID       int          `json:"user_id"`
	Type         FileType     `json:"file_type"`
	UploadStatus UploadStatus `json:"upload_status"`

	Bucket    string `json:"bucket"`
	FileName  string `json:"file_name"`
	ObjectKey string `json:"object_key"`
	PublicURL string `json:"public_url"`

	CreatedAt time.Time  `pg:",notnull,use_zero" json:"created_at"`
	UpdatedAt time.Time  `pg:",notnull,use_zero" json:"updated_at"`
	DeletedAt *time.Time `pg:",soft_delete" json:"-"`
}

// BeforeInsert Before insert trigger
func (o *File) BeforeInsert(c context.Context) (context.Context, error) {
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()

	return c, nil
}

// BeforeUpdate Before Update trigger
func (o *File) BeforeUpdate(c context.Context) (context.Context, error) {
	o.UpdatedAt = time.Now()

	return c, nil
}
