package authentication

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	models "github.com/kalaiselvan07/todo/authentication/models"
	"github.com/kalaiselvan07/todo/authentication/token"
)

func Login(context *gin.Context) {
	var loginObj models.LoginRequest
	if err := context.ShouldBindJSON(&loginObj); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
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

	var tokenCreationTime = time.Now().UTC()
	var expirationTime = tokenCreationTime.Add(time.Duration(5) * time.Minute)
	tokenString, err := token.GenrateToken(claims, expirationTime)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Token not created"})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{"token": tokenString})
}

func ValidateToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.Request.Header.Get("apikey")

		valid, claims := token.VerifyToken(tokenString)

		if !valid {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request"})
			context.Abort()
			return
		}

		if len(context.Keys) == 0 {
			context.Keys = make(map[string]interface{})
		}
		context.Keys["CompanyId"] = claims.CompanyId
		context.Keys["Username"] = claims.Username
		context.Keys["Role"] = claims.Role
	}
}
