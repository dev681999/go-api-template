package transport

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func fields(e *zerolog.Event, c echo.Context, req *http.Request, res *echo.Response, start time.Time) *zerolog.Event {
	return e.Str("remote_ip", c.RealIP()).
		Str("time", time.Since(start).String()).
		Str("host", req.Host).
		Str("request", fmt.Sprintf("%s %s", req.Method, req.RequestURI)).
		Int("status", res.Status).
		Int64("size", res.Size).
		Str("user_agent", req.UserAgent())
}

// logger is a middleware and zap to provide an "access log" like logging for each request.
func loggingMiddleware(log zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			defaultFields := func(e *zerolog.Event) *zerolog.Event {
				return fields(e, c, req, res, start)
			}
			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
				defaultFields = func(e *zerolog.Event) *zerolog.Event {
					return fields(e, c, req, res, start).Str("request_id", id)
				}
			}

			n := res.Status
			switch {
			case n >= 500:
				defaultFields(log.Error()).Msg("Server error")
			case n >= 400:
				defaultFields(log.Warn()).Msg("Client error")
			case n >= 300:
				defaultFields(log.Info()).Msg("Redirection")
			default:
				defaultFields(log.Info()).Msg("Success")
			}

			return nil
		}
	}
}
