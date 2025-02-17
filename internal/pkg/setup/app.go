package setup

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

type (
	App struct {
		fiberApp *fiber.App
	}
)

const (
	appName     = "pack-management"
	defaultPort = "8080"
)

func NewApp() *App {
	fiberApp := fiber.New(fiber.Config{
		AppName:                  appName,
		JSONEncoder:              json.Marshal,
		JSONDecoder:              json.Unmarshal,
		DisableStartupMessage:    false,
		EnablePrintRoutes:        false,
		EnableSplittingOnParsers: true,
	})

	return &App{
		fiberApp: fiberApp,
	}
}

func (a *App) Start(port string) {
	if port == "" {
		port = defaultPort
	}

	err := a.fiberApp.Listen(":" + port)
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}

		log.Println("Stopped serving new connections.")
	}
}

func (a *App) FiberApp() *fiber.App {
	return a.fiberApp
}

func (a *App) Shutdown() {
	shutdownC := make(chan os.Signal, 1)
	signal.Notify(shutdownC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownC

	log.Println("Shutting down...")

	defer func() {
		ctx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownRelease()

		if err := a.fiberApp.ShutdownWithContext(ctx); err != nil {
			log.Fatalf("HTTP shutdown error: %v", err)
		}

		log.Println("Shutdown complete.")
	}()
}
