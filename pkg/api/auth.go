package api

import (
	"github.com/gin-gonic/gin"
	"github.com/goodfoodcesi/tracking-api/pkg/config"
	"github.com/goodfoodcesi/tracking-api/pkg/token"
	"net/http"
)

func JWTInterceptor(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := token.ExtractTokenFromHeader(c.Request.Header.Get("Authorization"))
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		err := token.ValidateToken(tokenString, cfg)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
