package transport

import (
	"go-api-template/internal/openapi"
	"go-api-template/internal/user"
)

type server struct {
	user.Transport
}

// New returns a new OpenAPI Echo Server implementation
func New(userTransport user.Transport) openapi.ServerInterface {
	return &server{
		userTransport,
	}
}
