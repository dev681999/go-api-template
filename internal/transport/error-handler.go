package transport

import (
	apperr "go-api-template/internal/error"
	"go-api-template/internal/openapi"
	"net/http"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func errorHandler(logger zerolog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		logger.Debug().Str("errType", reflect.TypeOf(err).String()).Msg("")
		aerr, ok := err.(*apperr.Error)
		if !ok {
			if he, ok := err.(*echo.HTTPError); ok {
				if sre, ok := he.Internal.(*openapi3filter.SecurityRequirementsError); ok {
					if saerr, ok := sre.Errors[0].(*apperr.Error); ok && saerr != nil {
						aerr = apperr.New("transport", saerr.Reason, he.Code, nil)
					} else {
						aerr = apperr.New("transport", sre.Errors[0].Error(), he.Code, nil)
					}
				} else {
					msg, ok := he.Message.(string)
					if !ok {
						msg = he.Error()
					}

					aerr = apperr.New("transport", msg, he.Code, he.Internal)
				}
			} else {
				aerr = apperr.New("transport", err.Error(), http.StatusInternalServerError, err)
			}
		}
		logger.Error().Str("err", err.Error()).Msg("error handler")

		code := aerr.Code
		message := openapi.Error{
			Error:   aerr.String(),
			Message: aerr.Reason,
		}

		// Send response
		if !c.Response().Committed {
			if c.Request().Method == http.MethodHead { // Issue #608
				err = c.NoContent(aerr.Code)
			} else {
				err = c.JSON(code, message)
			}
			if err != nil {
				logger.Err(err).Msg("Error Handler")
				// e.Logger.Error(err)
			}
		}
	}
}
