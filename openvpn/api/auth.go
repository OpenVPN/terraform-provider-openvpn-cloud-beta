package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthConfig struct {
	Host         string
	ClientID     string
	ClientSecret string
}

type AuthData struct {
	AccessToken string `json:"access_token,omitempty"`
}

const TokenEndpoint = "/oauth/token"

func (c *Client) Authenticate(ctx context.Context) error {
	authData, err := authenticationRequest(ctx, c.client, c.authConfig)
	if err != nil {
		return err
	}
	c.authData = authData
	return nil
}

func authenticationRequest(ctx context.Context, httpClient HttpClient, authConfig *AuthConfig) (*AuthData, error) {
	request, err := createAuthRequest(authConfig)
	if err != nil {
		return nil, err
	}

	response, err := httpClient.Do(request.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	authResponse := &AuthData{}
	err = processJsonResponse(response, authResponse)
	if err != nil {
		return nil, err
	}

	return authResponse, nil
}

func createAuthRequest(authConfig *AuthConfig) (*http.Request, error) {
	requestBody := map[string]string{
		"grant_type": "client_credentials",
		"scope":      "default",
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", authConfig.apiUrl(TokenEndpoint), bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(authConfig.ClientID, authConfig.ClientSecret)
	return request, nil
}

func (d AuthData) AuthorizeRequest(request *http.Request) {
	request.Header.Set("Authorization", "Bearer "+d.AccessToken)
}

func (c AuthConfig) apiUrl(format string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s", c.Host, "/api/beta", fmt.Sprintf(format, a...))
}
