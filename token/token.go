package token

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	models "github.com/kalaiselvan07/todo/models"
)

const (
	jWTPrivateToken = "secretTokenSecretToken"
	ip              = "localhost"
)

func GenrateToken(claims *models.JwtClaims, expirationTime time.Time) (string, error) {
	claims.ExpiresAt = expirationTime.Unix()
	claims.IssuedAt = time.Now().UTC().Unix()
	claims.Issuer = ip

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jWTPrivateToken))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
