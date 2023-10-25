package main

import (
	"net/http"
	"time"

	models "authentication/models"
	"authentication/token"

	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {

	var loginObj models.LoginRequest

	err := c.ShouldBindJSON(&loginObj)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	var claims = &models.JwtClaims{}

	claims.CompanyId = "companyId"

	claims.Username = loginObj.UserName

	if loginObj.UserName == "administrator" {
		claims.Role = "admin"
	} else {
		claims.Role = "user"
	}

	tokenCreationTime := time.Now().UTC()
	expirationTime := tokenCreationTime.Add(5 * time.Minute)
	tokenString, err := token.GenerateToken(claims, expirationTime)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token not created"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}

func validateToken(context *gin.Context) {
	tokenString := context.Request.Header.Get("apikey")
	valid, claims := token.VerifyToken(tokenString)
	if !valid {
		context.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request"})
		context.Abort()
		return
	}
	context.IndentedJSON(http.StatusAccepted, claims)
}

func main() {
	router := gin.Default()
	router.POST("/authenticate", login)
	router.POST("/validate", validateToken)
	router.Run(":8088")
}
