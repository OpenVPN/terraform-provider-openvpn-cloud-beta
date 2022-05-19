package openvpn

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"terraform-provider-openvpn/openvpn/api"
	"testing"
	"time"
)

func getAuthenticatedClient(t *testing.T) *api.Client {
	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set", resource.TestEnvVar))
	}

	client := getTestApiClient()
	err := client.Authenticate(context.Background())
	require.NoError(t, err)
	return client
}

func getTestApiClient() *api.Client {
	authConfig := &api.AuthConfig{
		Host:         os.Getenv("OVPN_HOST"),
		ClientID:     os.Getenv("OVPN_CLIENT_ID"),
		ClientSecret: os.Getenv("OVPN_CLIENT_SECRET"),
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	client := api.NewClient(httpClient, authConfig)
	return client
}

func createTestHost(t *testing.T, client *api.Client, region string) *api.Host {
	hostName := "test-" + RandomString(8)
	host, err := client.CreateHost(context.Background(), &api.CreateHostRequest{
		Name:           hostName,
		Description:    t.Name() + " test description",
		Domain:         hostName + ".example.com",
		InternetAccess: "LOCAL",
		Connectors: []api.CreateConnectorRequest{
			{
				Name:        hostName,
				Description: "none",
				VpnRegionId: region,
			},
		},
	})
	require.NoError(t, err)
	return host
}

func deleteTestHost(t *testing.T, client *api.Client, hostID string) {
	err := client.DeleteHost(context.Background(), hostID)
	assert.NoError(t, err)
}

func createTestConnector(t *testing.T, client *api.Client, hostID string, regionId string) *api.Connector {
	connector, err := client.CreateConnector(context.Background(), &api.CreateConnectorData{
		Name:            "test-" + RandomString(8),
		Description:     t.Name() + " connector description",
		VpnRegionId:     regionId,
		NetworkItemId:   hostID,
		NetworkItemType: api.NetworkItemTypeHost,
	})
	require.NoError(t, err)
	return connector
}

func getDefaultRegionID(t *testing.T, client *api.Client) string {
	regions, err := client.ListRegions(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(regions), 1)
	regionId := regions[0].ID
	return regionId
}
