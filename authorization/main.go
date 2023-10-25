package main

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

func AuthorizationMiddleware(context *gin.Context) {
	validRole := context.Param("role")
	var claims JwtClaims
	if err := context.ShouldBindJSON(&claims); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad Request"})
		context.Abort()
		return
	}
	if claims.Role != validRole {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request"})
		context.Abort()
		return
	}
	context.IndentedJSON(http.StatusAccepted, nil)
}

func main() {
	router := gin.Default()
	router.POST("/authorize/:role", AuthorizationMiddleware)
	router.Run(":8089")
}
