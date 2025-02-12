package nagerdateapi

import (
	"context"
	"fmt"
	"net/http"
	"pack-management/internal/pkg/http/client"
)

type (
	holidayAPIClient struct {
		baseURL    string
		baseClient client.Client
	}
)

func NewHolidayAPIClient(baseClient client.Client, baseURL string) Client {
	return &holidayAPIClient{
		baseClient: baseClient,
		baseURL:    baseURL,
	}
}

func (c *holidayAPIClient) GetHolidays(
	ctx context.Context,
	countryCode string,
	year string,
) (HolidayResponse, error) {
	holidayResponse := HolidayResponse{}

	err := c.baseClient.Do(
		ctx,
		client.Request{
			Method: http.MethodGet,
			URL:    fmt.Sprintf("%s/PublicHolidays/%s/%s", c.baseURL, year, countryCode),
		},
		&holidayResponse,
	)
	if err != nil {
		return nil, err
	}

	return holidayResponse, nil
}
