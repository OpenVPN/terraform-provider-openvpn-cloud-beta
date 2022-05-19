package api

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestClient_GetHost(t *testing.T) {
	mockHttpClient := newMockHttpClient()
	authConfig := getAuthConfigTestData()
	ctx := context.Background()
	authData := &AuthData{AccessToken: "AccessToken"}

	t.Run("non-authenticated", func(t *testing.T) {
		client := NewClient(mockHttpClient, authConfig)
		_, err := client.ListRegions(ctx)
		assert.Error(t, err)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("authenticated", func(t *testing.T) {
		client := NewClient(mockHttpClient, authConfig)
		client.authData = authData

		expectedHost := &Host{
			ID: "7283bf33-0bda-4757-8983-470fad295763",
		}

		mockHttpClient.mockDo(t, expectedHost, func(request *http.Request) {
			assertRequestAuthorizedWithToken(t, request, authData.AccessToken)
			assert.True(t, strings.HasSuffix(request.URL.Path, expectedHost.ID))
		})

		host, err := client.GetHost(ctx, expectedHost.ID)

		assert.NoError(t, err)
		assert.Equal(t, expectedHost, host)
		mockHttpClient.AssertExpectations(t)
	})
}

func TestClient_CreateHost(t *testing.T) {
	mockHttpClient := newMockHttpClient()
	authConfig := getAuthConfigTestData()
	ctx := context.Background()
	authData := &AuthData{AccessToken: "AccessToken"}

	createHostRequest := &CreateHostRequest{
		Name:           "created-host",
		Description:    "created host description",
		Domain:         "test.test.com",
		InternetAccess: "BLOCKED",
		Connectors: []CreateConnectorRequest{
			{
				Name:        "test",
				Description: "info",
				VpnRegionId: "default",
			},
		},
	}

	t.Run("non-authenticated", func(t *testing.T) {
		client := NewClient(mockHttpClient, authConfig)
		_, err := client.CreateHost(ctx, createHostRequest)
		assert.Error(t, err)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("authenticated", func(t *testing.T) {
		client := NewClient(mockHttpClient, authConfig)
		client.authData = authData

		expectedHost := &Host{
			Name: createHostRequest.Name,
		}
		mockHttpClient.mockDo(t, expectedHost, func(request *http.Request) {
			assertRequestAuthorizedWithToken(t, request, authData.AccessToken)
			assert.True(t, strings.HasSuffix(request.URL.Path, expectedHost.ID))
		})

		host, err := client.CreateHost(ctx, createHostRequest)

		assert.NoError(t, err)
		assert.Equal(t, expectedHost, host)
		mockHttpClient.AssertExpectations(t)
	})
}

func TestClient_UpdateHost(t *testing.T) {
	mockHttpClient := newMockHttpClient()
	authConfig := getAuthConfigTestData()
	ctx := context.Background()
	authData := &AuthData{AccessToken: "AccessToken"}

	hostID := "7283bf33-0bda-4757-8983-470fad295763"

	updateHostRequest := &UpdateHostRequest{
		Name:           "created-host",
		Description:    "created host description",
		Domain:         "test.test.com",
		InternetAccess: "BLOCKED",
	}

	t.Run("non-authenticated", func(t *testing.T) {
		client := NewClient(mockHttpClient, authConfig)
		_, err := client.UpdateHost(ctx, hostID, updateHostRequest)
		assert.Error(t, err)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("authenticated", func(t *testing.T) {
		client := NewClient(mockHttpClient, authConfig)
		client.authData = authData

		expectedHost := &Host{
			ID:   hostID,
			Name: updateHostRequest.Name,
		}
		mockHttpClient.mockDo(t, expectedHost, func(request *http.Request) {
			assertRequestAuthorizedWithToken(t, request, authData.AccessToken)
			assert.True(t, strings.HasSuffix(request.URL.Path, expectedHost.ID))
		})

		host, err := client.UpdateHost(ctx, hostID, updateHostRequest)

		assert.NoError(t, err)
		assert.Equal(t, expectedHost, host)
		mockHttpClient.AssertExpectations(t)
	})
}

func TestClient_CreateDeleteHost_Real(t *testing.T) {
	skipIfNotAcceptance(t)
	authConfig, err := getAuthConfig()
	assert.NoError(t, err)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	client := NewClient(httpClient, authConfig)
	err = client.Authenticate(context.Background())
	require.NoError(t, err)

	regions, err := client.ListRegions(context.Background())
	require.NoError(t, err)

	createHostRequest := &CreateHostRequest{
		Name:           "created-host",
		Description:    "created host description",
		Domain:         "test.test.com",
		InternetAccess: "BLOCKED",
		Connectors: []CreateConnectorRequest{
			{
				Name:        "created-host-conn",
				Description: "created-host-conn description",
				VpnRegionId: regions[0].ID,
			},
		},
	}

	host, err := client.CreateHost(context.Background(), createHostRequest)
	assert.NoError(t, err)
	require.NotNil(t, host)
	assert.NotEmpty(t, host)
	assert.Equal(t, createHostRequest.Name, host.Name)
	assert.NotEmpty(t, host.Connectors)

	err = client.DeleteHost(context.Background(), host.ID)
	assert.NoError(t, err)
}

func TestClient_CreateUpdateDeleteHost_Real(t *testing.T) {
	skipIfNotAcceptance(t)
	authConfig, err := getAuthConfig()
	assert.NoError(t, err)

	ctx := context.Background()
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	client := NewClient(httpClient, authConfig)
	err = client.Authenticate(ctx)
	require.NoError(t, err)

	regions, err := client.ListRegions(context.Background())
	require.NoError(t, err)

	createHostRequest := &CreateHostRequest{
		Name:           "cud-host",
		Description:    "created host description",
		Domain:         "test.test.com",
		InternetAccess: "BLOCKED",
		Connectors: []CreateConnectorRequest{
			{
				Name:        "created-host-conn",
				Description: "created-host-conn description",
				VpnRegionId: regions[0].ID,
			},
		},
	}

	host, err := client.CreateHost(ctx, createHostRequest)
	assert.NoError(t, err)
	require.NotNil(t, host)
	assert.NotEmpty(t, host)
	assert.Equal(t, createHostRequest.Name, host.Name)
	assert.Equal(t, createHostRequest.Description, host.Description)

	updateHostRequest := &UpdateHostRequest{
		Name:           "cud-host-upd",
		Description:    "created host updated description",
		Domain:         createHostRequest.Domain,
		InternetAccess: createHostRequest.InternetAccess,
	}

	updatedHost, err := client.UpdateHost(ctx, host.ID, updateHostRequest)
	assert.NoError(t, err)
	require.NotNil(t, updatedHost)
	assert.NotEmpty(t, updatedHost)
	assert.Equal(t, updateHostRequest.Name, updatedHost.Name)
	assert.Equal(t, updateHostRequest.Description, updatedHost.Description)
	assert.Equal(t, host.Connectors, updatedHost.Connectors)

	err = client.DeleteHost(ctx, host.ID)
	assert.NoError(t, err)
}

func TestClient_CreateGetDeleteHost_Real(t *testing.T) {
	skipIfNotAcceptance(t)
	authConfig, err := getAuthConfig()
	assert.NoError(t, err)

	ctx := context.Background()
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}
	client := NewClient(httpClient, authConfig)
	err = client.Authenticate(ctx)
	require.NoError(t, err)

	regions, err := client.ListRegions(context.Background())
	require.NoError(t, err)

	createHostRequest := &CreateHostRequest{
		Name:           "cud-host",
		Description:    "created host description",
		Domain:         "test.test.com",
		InternetAccess: "BLOCKED",
		Connectors: []CreateConnectorRequest{
			{
				Name:        "created-host-conn",
				Description: "created-host-conn description",
				VpnRegionId: regions[0].ID,
			},
		},
	}

	host, err := client.CreateHost(ctx, createHostRequest)
	assert.NoError(t, err)
	require.NotNil(t, host)
	assert.NotEmpty(t, host)
	assert.Equal(t, createHostRequest.Name, host.Name)
	assert.Equal(t, createHostRequest.Description, host.Description)

	receivedHost, err := client.GetHost(ctx, host.ID)
	assert.NoError(t, err)
	require.NotNil(t, receivedHost)
	assert.NotEmpty(t, receivedHost)
	assert.Equal(t, host.ID, receivedHost.ID)
	assert.Equal(t, host.Connectors, receivedHost.Connectors)

	err = client.DeleteHost(ctx, host.ID)
	assert.NoError(t, err)
}
