package openvpn

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"regexp"
	"terraform-provider-openvpn/openvpn/api"
	"testing"
)

func TestResourceConnector_basic(t *testing.T) {
	t.Skip("Invalid test case")
	resourceName := "openvpn_connector.test"
	connectorName := "con_basic_" + RandomString(7)

	client := getAuthenticatedClient(t)

	regionId := getDefaultRegionID(t, client)
	host := createTestHost(t, client, regionId)

	t.Cleanup(func() {
		deleteTestHost(t, client, host.ID)
	})

	networkItemId := host.ID
	networkItemType := api.NetworkItemTypeHost

	request := &api.CreateConnectorData{
		Name:        connectorName,
		Description: connectorName + "Description",
		VpnRegionId: regionId,
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: defaultProviderFactory,
		CheckDestroy:      testAccCheckConnectorDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: resourceConnectorOutputConfig("test", networkItemId, networkItemType, request),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", request.Name),
					resource.TestCheckResourceAttr(resourceName, "description", request.Description),
					resource.TestCheckResourceAttr(resourceName, "vpn_region_id", request.VpnRegionId),
					resource.TestCheckResourceAttr(resourceName, "network_item_id", networkItemId),
					resource.TestCheckResourceAttr(resourceName, "network_item_type", string(networkItemType)),
					resource.TestCheckResourceAttrSet(resourceName, "profile"),
				),
			},
		},
	})
}

func TestResourceConnector_invalidNetworkItemType(t *testing.T) {
	t.Skip("Invalid test case")
	connectorName := "con_basic_" + RandomString(7)
	networkItemId := "50d3bfed-0b5d-4060-91f1-1a6a61ee3aa9"
	networkItemType := api.NetworkItemType("INVALID_TYPE")

	request := &api.CreateConnectorData{
		Name:        connectorName,
		Description: connectorName + "Description",
		VpnRegionId: "us-dev-1",
	}

	client := getAuthenticatedClient(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: defaultProviderFactory,
		CheckDestroy:      testAccCheckConnectorDestroy(client),
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(".*invalid value for.*Possible values are: HOST, NETWORK"),
				Config:      resourceConnectorOutputConfig("test", networkItemId, networkItemType, request),
			},
		},
	})
}

func resourceConnectorOutputConfig(name, networkItemId string, networkItemType api.NetworkItemType, request *api.CreateConnectorData) string {
	return fmt.Sprintf(`
provider "openvpn" {}

resource "openvpn_connector" "%s" {
	name = "%s"
	description = "%s"
	vpn_region_id = "%s"
	network_item_id = "%s"
	network_item_type = "%s"
}
`, name, request.Name, request.Description, request.VpnRegionId, networkItemId, networkItemType)

}

func testAccCheckConnectorDestroy(client *api.Client) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openvpn_connector" {
				continue
			}

			connectorId := rs.Primary.ID
			networkItemId := rs.Primary.Attributes["network_item_id"]
			networkItemType := api.NetworkItemType(rs.Primary.Attributes["networkItemType"])

			_, err := client.GetConnector(context.Background(), connectorId)
			if err == nil {
				continue
			}
			err = client.DeleteConnector(context.Background(), networkItemId, networkItemType, connectorId)
			if err != nil {
				return err
			}
			return fmt.Errorf("connector still exists")
		}

		return nil
	}
}
