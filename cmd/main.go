package main

import (
	"context"
	"fmt"
	_ "github.com/jaam8/online_song_library/docs"
	"github.com/jaam8/online_song_library/internal/api"
	"github.com/jaam8/online_song_library/internal/config"
	"github.com/jaam8/online_song_library/internal/repository"
	"github.com/jaam8/online_song_library/internal/service"
	"github.com/jaam8/online_song_library/pkg/logger"
	"github.com/jaam8/online_song_library/pkg/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logg, _ := logger.New(cfg.LogLevel)

	db, err := postgres.New(cfg.Postgres)
	if err != nil {
		logg.Fatal("failed to connect to database", zap.Error(err))
	}
	logg.Info("connected to database", zap.String("host", cfg.Postgres.Host))

	r := repository.New(db, logg)
	s := service.New(r, logg, cfg.SwaggerUrl)
	h := api.New(s, logg)

	e := echo.New()
	e.Use(api.LoggingMiddleware(logg))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	e.GET("/api/v1/songs", h.GetAllSongsHandler)
	e.POST("/api/v1/songs", h.CreateSongHandler)
	e.GET("/api/v1/songs/:id", h.GetSongHandler)
	e.PUT("/api/v1/songs/:id", h.UpdateSongHandler)
	e.DELETE("/api/v1/songs/:id", h.DeleteSongHandler)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	go func() {
		logg.Info(fmt.Sprintf("server starting on port :%s", cfg.RestPort))
		if err = e.Start(":" + cfg.RestPort); err != nil {
			logg.Fatal("server error", zap.Error(err))
		}
		logg.Info(fmt.Sprintf("server run on port :%s", cfg.RestPort))
	}()

	<-ctx.Done()
	logg.Info("server stopped")
}
