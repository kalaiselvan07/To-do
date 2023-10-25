package token

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	models "authentication/models"
)

const (
	jWTPrivateToken = "secretTokenSecretToken"
	ip              = "localhost"
)

func VerifyToken(tokenString string) (bool, *models.JwtClaims) {
	claims := &models.JwtClaims{}
	token, err := getTokenFromString(tokenString, claims)
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return false, claims
	}

	if token.Valid {
		if e := claims.Valid(); e == nil {
			return true, claims
		} else {
			fmt.Println("Error in claims:", e) // Print the error from custom claims validation
		}
	}
	return false, claims
}

func getTokenFromString(tokenString string, claims *models.JwtClaims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jWTPrivateToken), nil
	})
}

func GenerateToken(claims *models.JwtClaims, expirationTime time.Time) (string, error) {
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
