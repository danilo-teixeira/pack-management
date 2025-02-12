package pack_test

import (
	"net/http"
	"os"
	"pack-management/internal/pkg/config"
	"pack-management/internal/pkg/database"
	"pack-management/internal/pkg/http/client"
	"pack-management/internal/pkg/http/dogapi"
	"pack-management/internal/pkg/http/nagerdateapi"
	"pack-management/internal/services/pack"
	"pack-management/internal/services/person"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/h2non/gock"
)

var (
	shutdownServer func()
	clientApp      func(req *http.Request) (*http.Response, error)

	dogApiURL       = "http://dogapidog:1000"
	negerDateAPIURL = "http://datenagerat:1000"
)

func beforeAll() {
	cfg, err := config.NewConfig()
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

	baseClient := client.NewClient()
	dogAPIClient := dogapi.NewDogAPIClient(baseClient, dogApiURL)
	nagerDateAPIClient := nagerdateapi.NewHolidayAPIClient(baseClient, negerDateAPIURL)

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
