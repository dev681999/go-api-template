package mapper

import (
	"go-api-template/internal/openapi"
	"go-api-template/internal/storage"
)

// MapOpenAPIUploadStatusToStorageUploadStatus maps openapi.UploadStatus to storage.UploadStatus
func MapOpenAPIUploadStatusToStorageUploadStatus(s openapi.UploadStatus) storage.UploadStatus {
	switch s {
	case openapi.UploadStatus_uploaded:
		return storage.UploadStatusUploaded
	case openapi.UploadStatus_failed:
		return storage.UploadStatusFailed
	}

	return storage.UploadStatusPending
}

// MapStorageUploadStatusToOpenAPIUploadStatus maps storage.UploadStatus to openapi.UploadStatus
func MapStorageUploadStatusToOpenAPIUploadStatus(s storage.UploadStatus) openapi.UploadStatus {
	switch s {
	case storage.UploadStatusPending:
		return openapi.UploadStatus_pending
	case storage.UploadStatusFailed:
		return openapi.UploadStatus_failed
	}

	return openapi.UploadStatus_pending
}
