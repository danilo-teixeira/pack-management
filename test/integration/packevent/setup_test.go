package packevent_test

import (
	"context"
	"net/http"
	"os"
	"pack-management/internal/domain/holiday"
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
	ctx := context.Background()
	bunDB, app, shutdown := helpers.Setup()
	shutdownServer = shutdown

	baseClient := client.NewClient()
	dogAPIClient := dogapi.NewDogAPIClient(baseClient, dogApiURL)
	nagerDateAPIClient := nagerdateapi.NewHolidayAPIClient(baseClient, negerDateAPIURL)

	holidayRepo := holiday.NewMysqlRepository(&holiday.RepositoryParams{
		DB: bunDB,
	})
	holidaySvc := holiday.NewService(&holiday.ServiceParams{
		Repo:   holidayRepo,
		Client: nagerDateAPIClient,
	})

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
		Repo:           packRepo,
		PersonService:  personSvc,
		DogAPIClient:   dogAPIClient,
		HolidayService: holidaySvc,
	})
	pack.NewHTPPHandler(&pack.HandlerParams{
		Service: packSvc,
		App:     app,
	})

	packeventRepo := packevent.NewMysqlRepository(&packevent.RepositoryParams{
		DB: bunDB,
	})
	packeventSvc := packevent.NewService(ctx, &packevent.ServiceParams{
		Repo:        packeventRepo,
		PackService: packSvc,
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
