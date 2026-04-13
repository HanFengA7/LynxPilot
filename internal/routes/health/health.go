package health

import (
	healthapi "github.com/LychApe/LynxPilot/internal/api/health"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(r gin.IRoutes, db *gorm.DB) {
	api := healthapi.New(db)
	r.GET("/healthz", api.Healthz)
}
