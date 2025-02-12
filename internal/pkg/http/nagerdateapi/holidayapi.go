package nagerdateapi

import "context"

type (
	holidayAPIClient struct {
		url string
	}
)

func NewHolidayAPIClient(url string) Client {
	return &holidayAPIClient{
		url: url,
	}
}

func (c *holidayAPIClient) GetHolidays(
	ctx context.Context,
	countryCode string,
	year string,
) (HolidayResponse, error) {
	// TODO: Implement this
	return nil, nil
}
