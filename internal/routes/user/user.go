package user

import (
	userapi "github.com/LychApe/LynxPilot/internal/api/user"
	authmiddleware "github.com/LychApe/LynxPilot/internal/middleware/auth"
	userservice "github.com/LychApe/LynxPilot/internal/service/user"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup, userService *userservice.Service) {
	api := userapi.New(userService)

	publicGroup := r.Group("")
	registerPublicRoutes(publicGroup, api)

	privateGroup := r.Group("")
	privateGroup.Use(authmiddleware.Required(userService))
	registerPrivateRoutes(privateGroup, api)
}

func registerPublicRoutes(r *gin.RouterGroup, api *userapi.API) {
	r.POST("/auth/login", api.Login)
}

func registerPrivateRoutes(r *gin.RouterGroup, api *userapi.API) {
	_ = r
	_ = api
}
