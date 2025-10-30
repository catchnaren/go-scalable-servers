package middleware

import (
	"net/http"

	"github.com/catchnaren/go-scalable-servers/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthorizationMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// ctx.Abort() //403
		// Read authorization token from request
		tokenString := ctx.GetHeader("Authorization")

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token not found",
			})
			ctx.Abort()
			return
		}

		// Remove bearer from token
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// validate token with salt
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Config.JWTSaltKey), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			ctx.Abort()
			return
		}
		
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			ctx.Set("email", claims["email"])
			ctx.Set("name", claims["name"])
		}
		
		ctx.Next()
	}
}