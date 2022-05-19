package api

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	OvpnHostKey         = "OVPN_HOST"
	OvpnClientIdKey     = "OVPN_CLIENT_ID"
	OvpnClientSecretKey = "OVPN_CLIENT_SECRET"
	OvpnConfigPathKey   = "OVPN_CONFIG_PATH"

	OvpnDefaultConfigPathTemplate = "%s/.openvpn/config.toml"
)

func TestClient_Authenticate(t *testing.T) {
	mockHttpClient := newMockHttpClient()
	authConfig, err := getAuthConfig()
	assert.NoError(t, err)

	client := NewClient(mockHttpClient, authConfig)
	ctx := context.Background()
	mockResponseBody := &AuthData{AccessToken: "AccessToken"}

	mockHttpClient.mockDo(t, mockResponseBody, func(request *http.Request) {
		assert.Equal(t, ctx, request.Context())
		assertRequestHasOAuthSecrets(t, request, authConfig)
	})

	err = client.Authenticate(ctx)
	assert.NoError(t, err)
	assert.True(t, client.IsAuthenticated())

	assert.Equal(t, mockResponseBody, client.authData)
}

func TestClient_Authenticate_Real(t *testing.T) {
	authConfig, err := getAuthConfig()
	require.NoError(t, err)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	client := NewClient(httpClient, authConfig)
	err = client.Authenticate(context.Background())
	assert.NoError(t, err)
}

func getAuthConfigTestData() *AuthConfig {
	authConfig := &AuthConfig{
		Host:         "https://test.openvpn.com",
		ClientID:     "ClientID",
		ClientSecret: "ClientSecret",
	}
	return authConfig
}

func getAuthConfig() (*AuthConfig, error) {
	if authVarsPresentInEnv() {
		return getAuthConfigFromEnv(), nil
	}

	return getAuthConfigFromFile()
}

func getAuthConfigFromEnv() *AuthConfig {
	return &AuthConfig{
		Host:         os.Getenv(OvpnHostKey),
		ClientID:     os.Getenv(OvpnClientIdKey),
		ClientSecret: os.Getenv(OvpnClientSecretKey),
	}
}

func getAuthConfigFromFile() (*AuthConfig, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	authConfig := &AuthConfig{}
	err = toml.Unmarshal(file, authConfig)
	if err != nil {
		return nil, err
	}

	return authConfig, nil
}

func getConfigFilePath() (string, error) {
	configPath := os.Getenv(OvpnConfigPathKey)
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		configPath = fmt.Sprintf(OvpnDefaultConfigPathTemplate, home)
	}

	return configPath, nil
}

func authVarsPresentInEnv() bool {
	return os.Getenv(OvpnHostKey) != "" &&
		os.Getenv(OvpnClientIdKey) != "" &&
		os.Getenv(OvpnClientSecretKey) != ""
}

func assertRequestHasOAuthSecrets(t *testing.T, request *http.Request, authConfig *AuthConfig) {
	clientID, clientSecret, ok := request.BasicAuth()
	a := assert.New(t)
	a.True(ok, "missing basic auth")
	a.Equal(authConfig.ClientID, clientID)
	a.Equal(authConfig.ClientSecret, clientSecret)
}

func getRequestFromArgs(t *testing.T, args mock.Arguments) *http.Request {
	request, ok := args.Get(0).(*http.Request)
	require.True(t, ok, "invalid argument, required to be request")
	return request
}
