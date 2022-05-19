package openvpn

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"terraform-provider-openvpn/openvpn/api"
	"testing"
)

func TestDataSourceHost(t *testing.T) {
	dataSourceName := "data.openvpn_host.test"

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
				Config: dataHostOutputConfig(host.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "%", "7"),
					resource.TestCheckResourceAttr(dataSourceName, "id", host.ID),
					testCheckHostValuesAreSetExistingHost(dataSourceName, host),
					resource.TestCheckResourceAttr(dataSourceName, "connectors.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "connectors.0.%", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "connectors.0.id", connector.ID),
					resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr(dataSourceName, "connectors.0.name", connector.Name),
						resource.TestCheckResourceAttr(dataSourceName, "connectors.0.vpn_region_id", connector.VpnRegionId),
					),
					resource.TestCheckResourceAttr(dataSourceName, "system_subnets.#", "2"),
				),
			},
		},
	})
}

func testCheckHostValuesAreSetExistingHost(dataSourceName string, host *api.Host) resource.TestCheckFunc {
	return testCheckHostValuesAreSet(dataSourceName, api.CreateHostRequest{
		Name:           host.Name,
		Description:    host.Description,
		Domain:         host.Domain,
		InternetAccess: host.InternetAccess,
	})
}

func dataHostOutputConfig(id string) string {
	return fmt.Sprintf(`
provider "openvpn" {}

data "openvpn_host" "test" {
	id = "%s"
}
`, id)
}
