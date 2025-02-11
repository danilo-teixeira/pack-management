package naegerdateapi

import "context"

type (
	Client interface {
		GetHolidays(ctx context.Context, countryCode string, year string) (HolidayResponse, error)
	}

	HolidayResponse []Holiday

	Holiday struct {
		Date        string   `json:"date"`
		LocalName   string   `json:"localName"`
		Name        string   `json:"name"`
		CountryCode string   `json:"countryCode"`
		Fixed       bool     `json:"fixed"`
		Global      bool     `json:"global"`
		Countries   []string `json:"counties"`
		LaunchYear  string   `json:"launchYear"`
		Types       []string `json:"types"`
	}
)
