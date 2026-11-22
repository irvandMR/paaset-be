package health

import (
	"net/http"
	"tjus-be/internal/pkg/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// Method untuk pasang rute health ke router Gin
func (h *HealthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.GetHealthHandler)
	rg.GET("/db", h.GetHealthHandlerDB)
}

func (h *HealthHandler) GetHealthHandler(c *gin.Context) {
	util.ResponseSuccess(c, "Server is running", nil)
}

func (h *HealthHandler) GetHealthHandlerDB(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.PingContext(c.Request.Context()) != nil {
		util.ResponseError(c, http.StatusInternalServerError, "Database is unreachable")
		return
	}
	util.ResponseSuccess(c, "Database is healthy", nil)
}
