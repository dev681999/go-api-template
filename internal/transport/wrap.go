package transport

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// CustomContext .
type CustomContext struct {
	echo.Context
	logger zerolog.Logger
}

// HandlerFunc is custom context handler func
type HandlerFunc func(c *CustomContext) error

// Wrap wraps a custom handler into echo handler
func Wrap(fn HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*CustomContext)
		return fn(cc)
	}
}
