package auth

import (
	"net/http"
	"strings"

	userservice "github.com/LychApe/LynxPilot/internal/service/user"
	"github.com/gin-gonic/gin"
)

const ContextUsernameKey = "current_username"

func Required(userService *userservice.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader("Authorization"))
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(header, bearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(header, bearerPrefix))
		username, err := userService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set(ContextUsernameKey, username)
		c.Next()
	}
}
