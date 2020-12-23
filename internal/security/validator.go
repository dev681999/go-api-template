package security

import (
	"context"
	"net/http"

	apperr "go-api-template/internal/error"

	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
)

type (
	jwtExtractor func(echo.Context) (string, error)
)

// Errors
var (
	ErrJWTMissing = apperr.New("transport", "missing or malformed jwt", http.StatusBadRequest, nil)
	ErrJWTInvalid = apperr.New("transport", "invalid or expired jwt", http.StatusUnauthorized, nil)
)

// Defaults
const (
	ContextKey  = "user"
	TokenHeader = echo.HeaderAuthorization
	AuthScheme  = "Bearer"
)

// jwtFromHeader returns a `jwtExtractor` that extracts token from the request header.
func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:], nil
		}
		return "", ErrJWTMissing
	}
}

// jwtFromQuery returns a `jwtExtractor` that extracts token from the query string.
func jwtFromQuery(param string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		token := c.QueryParam(param)
		if token == "" {
			return "", ErrJWTMissing
		}
		return token, nil
	}
}

// jwtFromParam returns a `jwtExtractor` that extracts token from the url param string.
func jwtFromParam(param string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		token := c.Param(param)
		if token == "" {
			return "", ErrJWTMissing
		}
		return token, nil
	}
}

// jwtFromCookie returns a `jwtExtractor` that extracts token from the named cookie.
func jwtFromCookie(name string) jwtExtractor {
	return func(c echo.Context) (string, error) {
		cookie, err := c.Cookie(name)
		if err != nil {
			return "", ErrJWTMissing
		}
		return cookie.Value, nil
	}
}

// ClaimsParser is claims parser
type ClaimsParser interface {
	ParseTokenWithClaims(ctx context.Context, auth string) (interface{}, bool, error)
}

// ValidationMiddleware returns a new ehco validator middleware for openapi
func ValidationMiddleware(swagger *openapi3.Swagger, claimsParser ClaimsParser) echo.MiddlewareFunc {
	validatorOptions := &oapimiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
				ec := oapimiddleware.GetEchoContext(ctx)
				if ec == nil {
					return echo.ErrBadRequest
				}

				extractor := jwtFromHeader(echo.HeaderAuthorization, AuthScheme)
				auth, err := extractor(ec)
				if err != nil {
					return err
				}

				claims, valid, err := claimsParser.ParseTokenWithClaims(ctx, auth)
				if err == nil && valid {
					ec.Set(ContextKey, claims)
					return nil
				}

				if aerr, ok := err.(*apperr.Error); ok && aerr != nil {
					return aerr
				}

				return ErrJWTInvalid.CloneWithInner(err)
			},
		},
	}

	return oapimiddleware.OapiRequestValidatorWithOptions(swagger, validatorOptions)
}
