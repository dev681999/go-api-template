package transport

import (
	"go-api-template/internal/openapi"
	"go-api-template/internal/storage"
	"go-api-template/internal/user"
)

type server struct {
	user.TransportUser
	storage.TransportStorage
}

// New returns a new OpenAPI Echo Server implementation
func New(
	userTransport user.TransportUser,
	storageTransport storage.TransportStorage,
) openapi.ServerInterface {
	return &server{
		userTransport,
		storageTransport,
	}
}
