package pack_test

import (
	"net/http"
	"os"
	"pack-management/internal/pkg/database"
	"pack-management/internal/pkg/http/dogapi"
	"pack-management/internal/pkg/http/nagerdateapi"
	"pack-management/internal/services/pack"
	"pack-management/internal/services/person"
	"testing"

	"github.com/caarlos0/env"
	"github.com/gofiber/fiber/v2"
	"github.com/h2non/gock"
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

var shutdownServer func()
var clientApp func(req *http.Request) (*http.Response, error)

func beforeAll() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		AppName: "test",
	})

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

	clientApp = func(req *http.Request) (*http.Response, error) {
		req.Header.Set("Content-Type", "application/json")
		return app.Test(req, -1)
	}

	shutdownServer = func() {
		db.Close()
		app.Shutdown()
	}
}

func AfterAll() {
	gock.Off()
	shutdownServer()
}

func TestMain(m *testing.M) {
	beforeAll()
	code := m.Run()
	AfterAll()

	os.Exit(code)
}
