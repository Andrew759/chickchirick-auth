package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenS := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenS, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		//TODO: возможно не потребуется
		//Сохранение данных в контекст Gin, чтобы получать доступ к ним в контроллерах
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("uuid", claims["sub"])
		}

		c.Next()
	}
}
