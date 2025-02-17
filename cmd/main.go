package main

import (
	"log"
	"pack-management/internal/domain/holiday"
	"pack-management/internal/domain/metric"
	"pack-management/internal/domain/pack"
	"pack-management/internal/domain/packevent"
	"pack-management/internal/domain/person"
	"pack-management/internal/pkg/config"
	"pack-management/internal/pkg/database"
	"pack-management/internal/pkg/http/client"
	"pack-management/internal/pkg/http/dogapi"
	"pack-management/internal/pkg/http/nagerdateapi"
	"pack-management/internal/pkg/setup"
)

const appPort = "3300"

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	db, err := database.NewDatabase(&database.Params{
		DBHost:         cfg.DBHost,
		DBPort:         cfg.DBPort,
		DBName:         cfg.DBName,
		DBUser:         cfg.DBUser,
		DBPassword:     cfg.DBPassword,
		ConnectionPool: database.WithPoolConfigHigh(),
	}).Connect()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	baseAPP := setup.NewApp()
	fiberAPP := baseAPP.FiberApp()

	metric.NewHTPPHandler(&metric.HandlerParams{
		App: fiberAPP,
		DB:  db,
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

	holidayRepo := holiday.NewMysqlRepository(&holiday.RepositoryParams{
		DB: db,
	})
	holidaySvc := holiday.NewService(&holiday.ServiceParams{
		Repo:   holidayRepo,
		Client: nagerDateAPIClient,
	})

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
		Repo:           packRepo,
		PersonService:  personSvc,
		DogAPIClient:   dogAPIClient,
		HolidayService: holidaySvc,
	})
	pack.NewHTPPHandler(&pack.HandlerParams{
		Service: packSvc,
		App:     fiberAPP,
	})

	packEventRepo := packevent.NewMysqlRepository(&packevent.RepositoryParams{
		DB: db,
	})
	packEventSvc := packevent.NewService(&packevent.ServiceParams{
		Repo:        packEventRepo,
		PackService: packSvc,
	})
	packevent.NewHTPPHandler(&packevent.HandlerParams{
		Service: packEventSvc,
		App:     fiberAPP,
	})

	go baseAPP.Start(appPort)

	baseAPP.Shutdown()
}
