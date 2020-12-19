package security

import (
	"context"
	"fmt"
	"go-api-template/internal/user"

	oapimiddleware "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
)

// ValidationMiddleware returns a new ehco validator middleware for openapi
func ValidationMiddleware(swagger *openapi3.Swagger, jwtKey string, getUserFunc func(ctx context.Context, id int) (*user.User, error)) echo.MiddlewareFunc {
	validatorOptions := &oapimiddleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(c context.Context, input *openapi3filter.AuthenticationInput) error {
				ec := oapimiddleware.GetEchoContext(c)
				if ec == nil {
					return echo.ErrBadRequest
				}

				extractor := jwtFromHeader(echo.HeaderAuthorization, AuthScheme)
				auth, err := extractor(ec)
				if err != nil {
					return err
				}

				claims := &JwtClaims{}

				token, err := jwt.ParseWithClaims(auth, claims, func(t *jwt.Token) (interface{}, error) {
					// Check the signing method
					if t.Method.Alg() != AlgorithmHS256 {
						return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
					}

					return []byte(jwtKey), nil
				})
				if err == nil && token.Valid {
					_, err := getUserFunc(c, claims.UserID)
					if err == nil {
						// claims.UserID = u.ID
						// claims.UserID = u
						// Store user information from token into context.
						ec.Set(ContextKey, claims)
						return nil
					}
				}

				return &echo.HTTPError{
					Code:     ErrJWTInvalid.Code,
					Message:  ErrJWTInvalid.Message,
					Internal: err,
				}
			},
		},
	}

	return oapimiddleware.OapiRequestValidatorWithOptions(swagger, validatorOptions)
}
