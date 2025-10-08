package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Middleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error" : "authorization header missing", 
			})
			ctx.Abort()
			return 
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer")
		claims, err := ValidateToken(tokenString)
		if err != nil{
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token", 
			})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Next()
	}
}