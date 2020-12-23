package openapi

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// NewPrefixEchoMiddleware returns a new prefix echo middleware
func NewPrefixEchoMiddleware(prefix string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := c.Request()
			request.URL.Host += prefix
			request.URL.Path = strings.ReplaceAll(request.URL.Path, prefix, "")

			return next(c)
		}
	}
}
