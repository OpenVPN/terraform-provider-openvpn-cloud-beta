package api

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestClient_ListRegions(t *testing.T) {
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

		expectedRegions := []Region{
			{
				ID: "us-west-2",
			},
		}
		mockHttpClient.mockDo(t, expectedRegions, func(request *http.Request) {
			assertRequestAuthorizedWithToken(t, request, authData.AccessToken)
		})

		regions, err := client.ListRegions(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedRegions, regions)
		mockHttpClient.AssertExpectations(t)
	})
}

func TestClient_ListRegions_Real(t *testing.T) {
	authConfig, err := getAuthConfig()
	assert.NoError(t, err)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	client := NewClient(httpClient, authConfig)
	err = client.Authenticate(context.Background())
	require.NoError(t, err)

	regions, err := client.ListRegions(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, regions)
	assert.NotEmpty(t, regions)
	// t.Log(regions)
}

func assertRequestAuthorizedWithToken(t *testing.T, request *http.Request, accessToken string) bool {
	return assert.Equal(t, "Bearer "+accessToken, request.Header.Get("Authorization"))
}
