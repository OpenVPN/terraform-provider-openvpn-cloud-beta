package api

import (
	"context"
	"errors"
	"net/http"
)

type Region struct {
	ID         string `json:"id,omitempty"`
	Continent  string `json:"continent,omitempty"`
	Country    string `json:"country,omitempty"`
	CountryISO string `json:"countryIso,omitempty"`
	RegionName string `json:"regionName,omitempty"`
}

const RegionsEndpoint = "/regions"

func (c Client) ListRegions(ctx context.Context) ([]Region, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("authentication is required")
	}

	request, err := http.NewRequest("GET", c.apiEndpoint(RegionsEndpoint), nil)
	if err != nil {
		return nil, err
	}
	c.authData.AuthorizeRequest(request)

	response, err := c.client.Do(request.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	var regions []Region
	err = processJsonResponse(response, &regions)
	if err != nil {
		return nil, err
	}
	return regions, nil
}
