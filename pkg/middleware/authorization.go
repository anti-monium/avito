package middleware

import (
	api "avito_bootcamp/pkg/apartment_sale_api"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	token, _ := jwt.Parse(tokenString, api.ParseUserToken)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiresAt, err := claims.GetExpirationTime()
		if err != nil || time.Now().Unix() > expiresAt.Unix() {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		user_type, err := claims.GetSubject()
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Set("user_type", user_type)
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
