package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Authorization header is required"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Authorization header format must be Bearer {token}"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "default_secret"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Invalid or expired token"))
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
		} else {
			c.Error(utils.NewAppError(utils.ErrUnauthenticated, http.StatusUnauthorized, "Invalid token claims"))
			c.Abort()
			return
		}

		c.Next()
	}
}
