package models

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JwtClaims struct {
	CompanyId string `json:"companyId,omitempty"`
	Username  string `json:"username,omitempty"`
	Role      string `json:"role,omitempty"`
	jwt.StandardClaims
}

const ip = "localhost"

func (claims JwtClaims) Valid() error {
	var now = time.Now().UTC().Unix()
	if claims.VerifyExpiresAt(now, true) && claims.VerifyIssuer(ip, true) {
		return nil
	}
	return fmt.Errorf("token is invalid")
}
