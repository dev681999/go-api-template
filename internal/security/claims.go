package security

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// JwtClaims is a jwt claim
type JwtClaims struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserID    int    `json:"user_id"`
	jwt.StandardClaims
}

// GetClaimFromEchoContext gets users claims form context
func GetClaimFromEchoContext(c echo.Context) *JwtClaims {
	claims := c.Get("user").(*JwtClaims)

	return claims
}
