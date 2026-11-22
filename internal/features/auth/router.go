package auth

import (
	"time"
	"tjus-be/internal/config"
	"tjus-be/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetAuthRouter(rg *gin.RouterGroup, db *gorm.DB, env *config.Environment, rdb *redis.Client) {
	authRepository := NewAuthRepository(db)
	authService := NewAuthService(authRepository, env, rdb)
	authHandler := NewAuthHandler(authService)

	// Maksimal 5 percobaan login per menit per IP untuk mencegah Brute-Force
	loginLimiter := middleware.RateLimitMiddleware(rdb, "auth_login", 5, time.Minute)
	authHandler.RegisterRoutes(rg.Group("/auth"), loginLimiter)
}
