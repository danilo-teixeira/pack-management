package packevent_test

import (
	"net/http"
	"os"
	"pack-management/internal/domain/pack"
	"pack-management/internal/domain/packevent"
	"pack-management/internal/domain/person"
	"pack-management/internal/pkg/http/client"
	"pack-management/internal/pkg/http/dogapi"
	"pack-management/internal/pkg/http/nagerdateapi"
	"pack-management/test/helpers"
	"testing"
)

var (
	shutdownServer func()
	clientApp      func(req *http.Request) (*http.Response, error)

	dogApiURL       = "http://dogapidog:1000"
	negerDateAPIURL = "http://datenagerat:1000"
)

func beforeAll() {
	bunDB, app, shutdown := helpers.Setup()
	shutdownServer = shutdown

	baseClient := client.NewClient()
	dogAPIClient := dogapi.NewDogAPIClient(baseClient, dogApiURL)
	nagerDateAPIClient := nagerdateapi.NewHolidayAPIClient(baseClient, negerDateAPIURL)

	personRepo := person.NewMysqlRepository(&person.RepositoryParams{
		DB: bunDB,
	})
	personSvc := person.NewService(&person.ServiceParams{
		Repo: personRepo,
	})

	packRepo := pack.NewMysqlRepository(&pack.RepositoryParams{
		DB: bunDB,
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

	packeventRepo := packevent.NewMysqlRepository(&packevent.RepositoryParams{
		DB: bunDB,
	})
	packeventSvc := packevent.NewService(&packevent.ServiceParams{
		Repo: packeventRepo,
	})
	packevent.NewHTPPHandler(&packevent.HandlerParams{
		Service: packeventSvc,
		App:     app,
	})

	clientApp = func(req *http.Request) (*http.Response, error) {
		req.Header.Set("Content-Type", "application/json")
		return app.Test(req, -1)
	}
}

func AfterAll() {
	shutdownServer()
}

func TestMain(m *testing.M) {
	beforeAll()
	code := m.Run()
	AfterAll()

	os.Exit(code)
}
