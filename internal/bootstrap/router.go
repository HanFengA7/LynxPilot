package bootstrap

import (
	userrepo "github.com/LychApe/LynxPilot/internal/repository/user"
	healthroutes "github.com/LychApe/LynxPilot/internal/routes/health"
	userroutes "github.com/LychApe/LynxPilot/internal/routes/user"
	userservice "github.com/LychApe/LynxPilot/internal/service/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewGinEngine(db *gorm.DB, cfg *Config) *gin.Engine {
	r := gin.Default()

	healthroutes.Register(r, db)

	userRepository := userrepo.New(db)
	userService := userservice.New(userRepository, cfg.Auth.TokenSalt)
	api := r.Group("/api/v1")
	userroutes.Register(api, userService)

	return r
}
