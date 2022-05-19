package openvpn

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestDataSourceConnector(t *testing.T) {
	dataSourceName := "data.openvpn_connector.test"

	client := getAuthenticatedClient(t)

	regionId := getDefaultRegionID(t, client)

	host := createTestHost(t, client, regionId)

	t.Cleanup(func() {
		deleteTestHost(t, client, host.ID)
	})

	connector := host.Connectors[0]

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: defaultProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: dataConnectorOutputConfig(connector.ID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", connector.Name),
					resource.TestCheckResourceAttr(dataSourceName, "description", connector.Description),
					resource.TestCheckResourceAttr(dataSourceName, "ip_v4_address", connector.IpV4Address),
					resource.TestCheckResourceAttr(dataSourceName, "vpn_region_id", connector.VpnRegionId),
					resource.TestCheckResourceAttr(dataSourceName, "network_item_id", connector.NetworkItemId),
					resource.TestCheckResourceAttrSet(dataSourceName, "profile"),
				),
			},
		},
	})
}

func dataConnectorOutputConfig(connectorID string) string {
	return fmt.Sprintf(`
provider "openvpn" {
}

data "openvpn_connector" "test" {
	id = "%s"
}
`, connectorID)
}
