package health

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetHealthRouter(rg *gin.RouterGroup, db *gorm.DB) {
	healthHandler := NewHealthHandler(db)
	healthHandler.RegisterRoutes(rg.Group("/health"))
}
