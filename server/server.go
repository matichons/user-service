package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"user-service/config"
	"user-service/middlewares"
	_http "user-service/user/delivery/http"
	"user-service/user/repository"
	"user-service/user/usecase"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	httpServer  *http.Server
	cfg         *config.MainConfig
	db          *gorm.DB
	redisClient *redis.Client
}

func NewServer(cfg *config.MainConfig, db *gorm.DB, cache *redis.Client) *Server {
	return &Server{
		cfg:         cfg,
		db:          db,
		redisClient: cache,
	}
}
func (s *Server) Run() error {
	router := gin.Default()

	router.Use(middlewares.InitMiddleware().CORS())

	userRepo := repository.New(s.db)
	userCache := repository.NewCacheRepository(s.redisClient)

	timeoutContext := time.Duration(s.cfg.Connection.Timeout) * time.Second

	userUC := usecase.NewUserUsecase(userRepo, userCache, timeoutContext)
	_http.NewUserHandler(router, userUC)

	s.httpServer = &http.Server{
		Addr:           s.cfg.Server.Address,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return s.httpServer.Shutdown(ctx)
}
