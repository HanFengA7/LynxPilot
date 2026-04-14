package user

import (
	apiUser "github.com/LychApe/LynxPilot/internal/api/user"
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine) {
	userPublicGroup := router.Group("/api/v1/public/user")
	{
		//[登录]
		userPublicGroup.POST("/login", apiUser.LoginHandler)
	}
}
