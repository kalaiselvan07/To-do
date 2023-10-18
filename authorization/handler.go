package authorization

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authorization(validRole string) gin.HandlerFunc {
	return func(context *gin.Context) {
		if len(context.Keys) == 0 {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request"})
			context.Abort()
			return
		}
		rolesVal := context.Keys["Role"]
		if rolesVal == nil {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request"})
			context.Abort()
			return
		}

		if validRole != rolesVal {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized request"})
			context.Abort()
			return
		}
	}
}
