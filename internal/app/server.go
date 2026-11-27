package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tjus-be/internal/config"
	"tjus-be/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	Engine *gin.Engine
	Cfg    *config.Config
	DB     *gorm.DB
	Env    *config.Environment
	Redis  *redis.Client
}

func (s *Server) SetupConfiguration() {
	cfg := config.InitConfig()
	s.Cfg = cfg
	s.DB = cfg.DB
	s.Env = cfg.Env
	s.Redis = cfg.Redis
}

func (s *Server) SetupMiddlewares() {
	s.Engine.Use(gin.Recovery())
	s.Engine.Use(middleware.CORSMiddleware())
	s.Engine.Use(middleware.LoggerMiddleware(s.Env.IsProduction()))
}

func (s *Server) SetupRouters() {
	SetRouter(s)
}

func (s *Server) Close() {
	if s.Cfg != nil {
		s.Cfg.Close()
	}
}

func NewServer() *Server {
	server := &Server{}
	server.SetupConfiguration()

	if server.Env.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	server.Engine = gin.New()
	server.SetupMiddlewares()
	server.SetupRouters()

	return server
}

func (s *Server) Start() {
	srv := &http.Server{
		Addr:    ":" + s.Env.Port,
		Handler: s.Engine,
	}

	go func() {
		log.Printf("Server Paaset berjalan di port: %s\n", s.Env.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listen error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Mematikan server Paaset secara aman (Graceful Shutdown)...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server dipaksa berhenti: %v", err)
	}

	log.Println("Server berhasil dimatikan.")
}
