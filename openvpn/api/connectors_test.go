package api

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestClient_GetConnector(t *testing.T) {
	mockHttpClient := newMockHttpClient()
	authConfig := getAuthConfigTestData()
	ctx := context.Background()
	authData := &AuthData{AccessToken: "AccessToken"}

	t.Run("non-authenticated", func(t *testing.T) {
		// given
		client := NewClient(mockHttpClient, authConfig)

		// when
		_, err := client.GetConnector(ctx, "123")

		// then
		assert.Error(t, err)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("authenticated", func(t *testing.T) {
		// given
		client := NewClient(mockHttpClient, authConfig)
		client.authData = authData

		connectorId := "connector-id-123"
		expectedConnector := &Connector{
			ID:              connectorId,
			Name:            "returned Name",
			IpV4Address:     "returned IpV4Address",
			IpV6Address:     "returned IpV6Address",
			NetworkItemId:   "returned NetworkItemId",
			NetworkItemType: "returned NetworkItemType",
			VpnRegionId:     "returned VpnRegionId",
		}

		mockHttpClient.mockDo(t, expectedConnector, func(request *http.Request) {
			assertRequestAuthorizedWithToken(t, request, authData.AccessToken)
			assert.True(t, strings.HasSuffix(request.URL.Path, expectedConnector.ID))
		})

		// when
		connector, err := client.GetConnector(ctx, connectorId)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedConnector, connector)
		mockHttpClient.AssertExpectations(t)
	})
}

func TestClient_DeleteConnector(t *testing.T) {
	mockHttpClient := newMockHttpClient()
	authConfig := getAuthConfigTestData()
	ctx := context.Background()
	authData := &AuthData{AccessToken: "AccessToken"}

	t.Run("non-authenticated", func(t *testing.T) {
		// given
		client := NewClient(mockHttpClient, authConfig)

		// when
		err := client.DeleteConnector(ctx, "networkId123", NetworkItemTypeHost, "connectorID")

		// then
		assert.Error(t, err)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("authenticated", func(t *testing.T) {
		// given
		client := NewClient(mockHttpClient, authConfig)
		client.authData = authData

		connectorId := "connector-id-123"

		mockHttpClient.mockDo(t, nil, func(request *http.Request) {
			assertRequestAuthorizedWithToken(t, request, authData.AccessToken)

			reqUrl := request.URL
			assert.Equal(t, "HOST", reqUrl.Query().Get("networkItemType"))
			assert.Equal(t, "networkId123", reqUrl.Query().Get("networkItemId"))
			assert.True(t, strings.HasSuffix(reqUrl.Path, connectorId))
		})

		// when
		err := client.DeleteConnector(ctx, "networkId123", "HOST", connectorId)

		// then
		assert.NoError(t, err)
		mockHttpClient.AssertExpectations(t)
	})
}

func TestClient_CreateConnector(t *testing.T) {
	mockHttpClient := newMockHttpClient()
	authConfig := getAuthConfigTestData()
	ctx := context.Background()
	authData := &AuthData{AccessToken: "AccessToken"}

	t.Run("non-authenticated", func(t *testing.T) {
		// given
		client := NewClient(mockHttpClient, authConfig)

		createConnectorRequest := &CreateConnectorData{
			Name: "name123",
		}
		// when
		connector, err := client.CreateConnector(ctx, createConnectorRequest)

		// then
		assert.Error(t, err)
		assert.Nil(t, connector)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("authenticated", func(t *testing.T) {
		// given
		client := NewClient(mockHttpClient, authConfig)
		client.authData = authData

		createConnectorRequest := &CreateConnectorData{
			Name:            "name123",
			VpnRegionId:     "region",
			NetworkItemId:   "hostId78",
			NetworkItemType: NetworkItemTypeHost,
		}

		response := &Connector{
			ID:   "ID123",
			Name: "name123",
		}

		mockHttpClient.mockDo(t, response, func(request *http.Request) {
			assertRequestAuthorizedWithToken(t, request, authData.AccessToken)

			reqUrl := request.URL
			assert.Equal(t, string(createConnectorRequest.NetworkItemType), reqUrl.Query().Get("networkItemType"))
			assert.Equal(t, createConnectorRequest.NetworkItemId, reqUrl.Query().Get("networkItemId"))

		})

		// when
		connector, err := client.CreateConnector(ctx, createConnectorRequest)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, connector)
		assert.Equal(t, response, connector)
		mockHttpClient.AssertExpectations(t)
	})
}

func TestClient_UpdateConnector(t *testing.T) {
	mockHttpClient := newMockHttpClient()
	authConfig := getAuthConfigTestData()
	ctx := context.Background()
	authData := &AuthData{AccessToken: "AccessToken"}

	updateConnectorRequest := &CreateConnectorData{
		Name:            "name123",
		NetworkItemId:   "hostId78",
		NetworkItemType: NetworkItemTypeHost,
	}

	t.Run("non-authenticated", func(t *testing.T) {
		// given
		client := NewClient(mockHttpClient, authConfig)

		// when
		_, err := client.UpdateConnector(ctx, "123", updateConnectorRequest)

		// then
		assert.Error(t, err)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("authenticated", func(t *testing.T) {
		// given
		client := NewClient(mockHttpClient, authConfig)
		client.authData = authData

		connectorId := "connector-id-123"
		expectedConnector := &Connector{
			ID:              connectorId,
			Name:            "returned Name",
			IpV4Address:     "returned IpV4Address",
			IpV6Address:     "returned IpV6Address",
			NetworkItemId:   "returned NetworkItemId",
			NetworkItemType: "returned NetworkItemType",
			VpnRegionId:     "returned VpnRegionId",
		}

		mockHttpClient.mockDo(t, expectedConnector, func(request *http.Request) {
			assertRequestAuthorizedWithToken(t, request, authData.AccessToken)

			reqUrl := request.URL
			assert.Equal(t, string(updateConnectorRequest.NetworkItemType), reqUrl.Query().Get("networkItemType"))
			assert.Equal(t, updateConnectorRequest.NetworkItemId, reqUrl.Query().Get("networkItemId"))

			assert.True(t, strings.HasSuffix(reqUrl.Path, expectedConnector.ID))
		})

		// when
		connector, err := client.UpdateConnector(ctx, connectorId, updateConnectorRequest)

		// then
		assert.NoError(t, err)
		assert.Equal(t, expectedConnector, connector)
		mockHttpClient.AssertExpectations(t)
	})
}

func TestClient_GetConnectorProfile(t *testing.T) {
	mockHttpClient := newMockHttpClient()
	authConfig := getAuthConfigTestData()
	ctx := context.Background()
	authData := &AuthData{AccessToken: "AccessToken"}

	t.Run("non-authenticated", func(t *testing.T) {
		// given
		client := NewClient(mockHttpClient, authConfig)

		// when
		_, err := client.GetConnectorProfile(ctx, "123")

		// then
		assert.Error(t, err)
		mockHttpClient.AssertExpectations(t)
	})

	t.Run("authenticated", func(t *testing.T) {
		// given
		client := NewClient(mockHttpClient, authConfig)
		client.authData = authData

		connectorId := "connector-id-123"

		exampleConnectorProfile := "Profile text\nMultiline"
		mockHttpClient.mockDoBytes(t, []byte(exampleConnectorProfile), func(request *http.Request) {
			assertRequestAuthorizedWithToken(t, request, authData.AccessToken)
			assert.True(t, strings.Contains(request.URL.Path, connectorId))
		})

		// when
		connectorProfile, err := client.GetConnectorProfile(ctx, connectorId)

		// then
		assert.NoError(t, err)
		assert.Equal(t, exampleConnectorProfile, connectorProfile)
		mockHttpClient.AssertExpectations(t)
	})
}

func TestClient_UpdateConnector_safe_info_Real(t *testing.T) {
	// setup
	skipIfNotAcceptance(t)
	authConfig, err := getAuthConfig()
	assert.NoError(t, err)

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFunc()
	httpClient := &http.Client{
		Timeout: 35 * time.Second,
	}

	client := NewClient(httpClient, authConfig)
	err = client.Authenticate(ctx)
	require.NoError(t, err)

	regions, err := client.ListRegions(context.Background())
	require.NoError(t, err)

	// create host
	suffix := randomString()

	createHostRequest := &CreateHostRequest{
		Name:           "host name " + suffix,
		Description:    "created host description ",
		Domain:         suffix + ".domain.com",
		InternetAccess: "BLOCKED",
		Connectors: []CreateConnectorRequest{
			{
				Name:        "c_name_" + suffix,
				Description: "created connector description",
				VpnRegionId: regions[0].ID,
			},
		},
	}

	host, err := client.CreateHost(ctx, createHostRequest)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := client.DeleteHost(context.Background(), host.ID)
		assert.NoError(t, err)
	})

	require.NotEmpty(t, host.Connectors)

	createdConnector := host.Connectors[0]

	// create connector
	updateRequest := CreateConnectorData{
		Name:            createdConnector.Name + "_2",
		Description:     createdConnector.Description,
		VpnRegionId:     createdConnector.VpnRegionId,
		NetworkItemId:   host.ID,
		NetworkItemType: NetworkItemTypeHost,
	}

	updatedConnector, err := client.UpdateConnector(ctx, host.Connectors[0].ID, &updateRequest)
	require.NoError(t, err)

	assert.Equal(t, createdConnector.ID, updatedConnector.ID)
	assert.Equal(t, updateRequest.Name, updatedConnector.Name)
	assert.Equal(t, updateRequest.Description, updatedConnector.Description)
	assert.Equal(t, updateRequest.VpnRegionId, updatedConnector.VpnRegionId)
}

func TestClient_UpdateConnector_changing_region_Real(t *testing.T) {
	// setup
	skipIfNotAcceptance(t)
	authConfig, err := getAuthConfig()
	assert.NoError(t, err)

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFunc()
	httpClient := &http.Client{
		Timeout: 35 * time.Second,
	}

	client := NewClient(httpClient, authConfig)
	err = client.Authenticate(ctx)
	require.NoError(t, err)

	regions, err := client.ListRegions(context.Background())
	require.NoError(t, err)

	// create host
	suffix := randomString()

	createHostRequest := &CreateHostRequest{
		Name:           "host name " + suffix,
		Description:    "created host description ",
		Domain:         suffix + ".domain.com",
		InternetAccess: "BLOCKED",
		Connectors: []CreateConnectorRequest{
			{
				Name:        "c_name_" + suffix,
				Description: "created connector description",
				VpnRegionId: regions[0].ID,
			},
		},
	}

	host, err := client.CreateHost(ctx, createHostRequest)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := client.DeleteHost(context.Background(), host.ID)
		assert.NoError(t, err)
	})

	require.NotEmpty(t, host.Connectors)

	createdConnector := host.Connectors[0]

	// create connector
	updateRequest := CreateConnectorData{
		Name:            createdConnector.Name,
		Description:     createdConnector.Description,
		VpnRegionId:     regions[1].ID,
		NetworkItemId:   host.ID,
		NetworkItemType: NetworkItemTypeHost,
	}

	updatedConnector, err := client.UpdateConnector(ctx, host.Connectors[0].ID, &updateRequest)
	require.NoError(t, err)

	assert.Equal(t, createdConnector.ID, updatedConnector.ID)
	assert.Equal(t, updateRequest.Name, updatedConnector.Name)
	assert.Equal(t, updateRequest.Description, updatedConnector.Description)
	assert.Equal(t, updateRequest.VpnRegionId, updatedConnector.VpnRegionId)
}

func skipIfNotAcceptance(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Skipping test, use TF_ACC=1 to enable it")
	}
}

func randomString() string {
	randomString, err := uuid.GenerateUUID()
	if err != nil {
		return ""
	}
	return randomString[:7]
}
