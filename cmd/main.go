package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pack-management/internal/domain/pack"
	"pack-management/internal/domain/person"
	"pack-management/internal/pkg/config"
	"pack-management/internal/pkg/database"
	"pack-management/internal/pkg/http/client"
	"pack-management/internal/pkg/http/dogapi"
	"pack-management/internal/pkg/http/nagerdateapi"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	appName = "pack-management"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	db, err := database.NewDatabase(&database.Params{
		DBHost:     cfg.DBHost,
		DBPort:     cfg.DBPort,
		DBName:     cfg.DBName,
		DBUser:     cfg.DBUser,
		DBPassword: cfg.DBPassword,
	}).Connect()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// TODO: move to separate package???
	app := fiber.New(fiber.Config{
		AppName:                  appName,
		JSONEncoder:              json.Marshal,
		JSONDecoder:              json.Unmarshal,
		DisableStartupMessage:    true,
		EnablePrintRoutes:        false,
		EnableSplittingOnParsers: true,
	})

	baseClient := client.NewClient()
	dogAPIClient := dogapi.NewDogAPIClient(
		baseClient,
		"https://dogapi.dog/api/v2",
	)
	nagerDateAPIClient := nagerdateapi.NewHolidayAPIClient(
		baseClient,
		"https://date.nager.at/api/v3",
	)

	personRepo := person.NewMysqlRepository(&person.RepositoryParams{
		DB: db,
	})
	personSvc := person.NewService(&person.ServiceParams{
		Repo: personRepo,
	})

	packRepo := pack.NewMysqlRepository(&pack.RepositoryParams{
		DB: db,
	})
	packSvc := pack.NewService(&pack.ServiceParams{
		Repo:               packRepo,
		PersonService:      personSvc,
		DogAPIClient:       dogAPIClient,
		NagerDateAPIClient: nagerDateAPIClient,
	})
	pack.NewHTPPHandler(&pack.HandlerParams{
		Service: packSvc,
		App:     app,
	})

	go func() {
		err = app.Listen(":3300")
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("HTTP server error: %v", err)
			}

			log.Println("Stopped serving new connections.")
		}
	}()

	// TODO: move to separate package???
	shutdownC := make(chan os.Signal, 1)
	signal.Notify(shutdownC, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownC

	ctx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Shutdown complete.")
}
