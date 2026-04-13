package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type API struct {
	db *gorm.DB
}

func New(db *gorm.DB) *API {
	return &API{db: db}
}

func (a *API) Healthz(c *gin.Context) {
	sqlDB, err := a.db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "down", "error": err.Error()})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "down", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"database": "sqlite",
	})
}
