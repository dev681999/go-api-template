package transport

import (
	"io/ioutil"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

// NewEchoEngine returns a new echo server
func NewEchoEngine(logger zerolog.Logger) *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = errorHandler(logger)

	e.Logger.SetOutput(ioutil.Discard)

	/* e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{
				Context: c,
				logger:  logger,
			}
			return next(cc)
		}
	}) */
	e.Use(loggingMiddleware(logger))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return e
}
