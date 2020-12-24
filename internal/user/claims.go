package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// ClaimType is type of claim
type ClaimType uint

// ClaimTypes
const (
	ClaimTypeNormal ClaimType = iota
	ClaimTypeActivation
	ClaimTypePasswordReset
)

// Role is role of the user
type Role uint

// Roles
const (
	RoleAdmin Role = iota
	RoleUser
)

// Claims is a jwt claim
type Claims struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	UserID    int       `json:"user_id"`
	Role      Role      `json:"role"`
	Type      ClaimType `json:"claim_type"`
	Email     string    `json:"-"`
	Mobile    string    `json:"-"`

	jwt.StandardClaims
}

// GetClaimFromEchoContext gets users claims form context
func GetClaimFromEchoContext(c echo.Context) *Claims {
	claims := c.Get("user").(*Claims)

	return claims
}
