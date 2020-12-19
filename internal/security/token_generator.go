package security

import (
	"go-api-template/internal/user"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken returns a function to  genrate a jwt token for a user
func GenerateToken(jwtKey string) func(u *user.User) (string, error) {
	return func(u *user.User) (string, error) {
		// Set custom claims
		claims := &JwtClaims{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			UserID:    u.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24 * 365).Unix(),
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(jwtKey))
		if err != nil {
			return "", err
		}

		return t, nil
	}
}
