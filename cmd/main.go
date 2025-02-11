package main

import (
	"pack-management/internal/pkg/database"
	"pack-management/internal/pkg/http/dogapi"
	nagerdateapi "pack-management/internal/pkg/http/nagerdateapi"
	"pack-management/internal/services/pack"
	"pack-management/internal/services/person"

	"github.com/caarlos0/env"
	"github.com/gofiber/fiber/v2"
)

type (
	config struct {
		DBHost     string `env:"DB_HOST,required"`
		DBPort     string `env:"DB_PORT,required"`
		DBName     string `env:"DB_NAME,required"`
		DBUser     string `env:"DB_USER,required"`
		DBPassword string `env:"DB_PASSWORD,required"`
	}
)

func main() {
	// TODO: move to config package
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	db, err := database.NewDatabase(&database.Params{
		DBHost:     cfg.DBHost,
		DBPort:     cfg.DBPort,
		DBName:     cfg.DBName,
		DBUser:     cfg.DBUser,
		DBPassword: cfg.DBPassword,
	}).Connect()
	if err != nil {
		panic(err)
	}

	_, err = db.Query("SELECT 1")
	if err != nil {
		panic("Failed to connect to database")
	}

	// TODO: move to separate package???
	app := fiber.New()

	dogAPIClient := dogapi.NewDogAPIClient("https://dogapi.dog/api/v2")
	nagerDateAPIClient := nagerdateapi.NewHolidayAPIClient("https://date.nager.at/api/v3")

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

	err = app.Listen(":3300")
	if err != nil {
		panic(err)
	}
}
