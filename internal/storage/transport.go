package storage

import (
	apperr "go-api-template/internal/error"
	"go-api-template/internal/openapi"
	"go-api-template/internal/user"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// TransportStorage handles transport for service
type TransportStorage struct {
	logger zerolog.Logger
	srv    Service
}

// NewTransport creates a new transport
func NewTransport(
	logger zerolog.Logger,
	srv Service,
) TransportStorage {
	return TransportStorage{
		logger: logger,
		srv:    srv,
	}
}

// CreateFile creates a new file and returns a presigned url to upload file
func (t TransportStorage) CreateFile(c echo.Context) error {
	ctx := c.Request().Context()

	claims := user.GetClaimFromEchoContext(c)

	req := &openapi.CreateFileJSONBody{}
	err := c.Bind(req)
	if err != nil {
		t.logger.Err(err).Msg("")
		return apperr.NewWithCode("transport", "bad request", http.StatusBadRequest, err)
	}

	f, presignedURL, err := t.srv.Create(ctx, &File{
		UserID:   claims.UserID,
		FileName: req.Filename,
	})
	if err != nil {

	}

	return c.JSON(http.StatusOK, openapi.UploadFileResponse{
		Filename:     f.FileName,
		Id:           f.ID,
		PresignedUrl: presignedURL,
	})
}
