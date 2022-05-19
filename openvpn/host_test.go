package openvpn

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
	"terraform-provider-openvpn/openvpn/api"
	"testing"
)

func TestResourceHost_basic(t *testing.T) {
	resourceName := "openvpn_host.test"

	hostName := "test-" + RandomString(8)

	client := getAuthenticatedClient(t)

	regions, err := client.ListRegions(context.Background())
	require.NoError(t, err)

	hostValues := api.CreateHostRequest{
		Name:           hostName,
		Description:    hostName + " Description",
		Domain:         hostName + ".example.com",
		InternetAccess: "BLOCKED",
		Connectors: []api.CreateConnectorRequest{
			{
				Name:        hostName + "_c",
				Description: "Connector for host " + hostName,
				VpnRegionId: regions[0].ID,
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: defaultProviderFactory,
		CheckDestroy:      testAccCheckHostDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: resourceHostOutputConfig("test", hostValues),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testCheckHostValuesAreSet(resourceName, hostValues),
					resource.TestCheckResourceAttr(resourceName, "connector.#", "1"),
					testCheckHostConnectorValues(resourceName, hostValues.Connectors),
				),
			},
		},
	})
}

func TestResourceHost_edit(t *testing.T) {
	resourceName := "openvpn_host.test"
	hostName := "test-" + RandomString(7)
	client := getAuthenticatedClient(t)

	regions, err := client.ListRegions(context.Background())
	require.NoError(t, err)

	hostValues := api.CreateHostRequest{
		Name:           hostName,
		Description:    hostName + " Description",
		Domain:         hostName + ".example.com",
		InternetAccess: "BLOCKED",
		Connectors: []api.CreateConnectorRequest{
			{
				Name:        hostName + "_c",
				Description: "Connector for host " + hostName,
				VpnRegionId: regions[0].ID,
			},
		},
	}
	newHostValues := api.CreateHostRequest{
		Name:           hostName + "-2",
		Description:    hostName + "Description 2",
		Domain:         hostName + ".example.com",
		InternetAccess: "BLOCKED",
		Connectors: []api.CreateConnectorRequest{
			{
				Name:        hostName + "_e",
				Description: "Connector edited for host " + hostName,
				VpnRegionId: regions[1].ID,
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: defaultProviderFactory,
		CheckDestroy:      testAccCheckHostDestroy(client),
		Steps: []resource.TestStep{
			{
				Config: resourceHostOutputConfig("test", hostValues),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testCheckHostValuesAreSet(resourceName, hostValues),
					testCheckHostConnectorValues(resourceName, hostValues.Connectors),
				),
			},
			{
				Config: resourceHostOutputConfig("test", newHostValues),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testCheckHostValuesAreSet(resourceName, newHostValues),
					testCheckHostConnectorValues(resourceName, newHostValues.Connectors),
				),
			},
		},
	})
}

func testAccCheckHostDestroy(client *api.Client) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openvpn_host" {
				continue
			}

			hostID := rs.Primary.ID

			_, err := client.GetHost(context.Background(), hostID)
			if err == nil {
				err := client.DeleteHost(context.Background(), hostID)
				if err != nil {
					return nil
				}
				return fmt.Errorf("still exists")
			}
		}

		return nil
	}
}

func testCheckHostValuesAreSet(dataSourceName string, host api.CreateHostRequest) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(dataSourceName, "name", host.Name),
		resource.TestCheckResourceAttr(dataSourceName, "description", host.Description),
		resource.TestCheckResourceAttr(dataSourceName, "internet_access", host.InternetAccess),
		resource.TestCheckResourceAttr(dataSourceName, "domain", host.Domain),
	)
}

func testCheckHostConnectorValues(resourceName string, connectorRequests []api.CreateConnectorRequest) resource.TestCheckFunc {
	tests := make([]resource.TestCheckFunc, len(connectorRequests))
	for i, connectorRequest := range connectorRequests {
		tests[i] = testCheckHostConnectorItemValues(resourceName, i, connectorRequest)
	}
	return resource.ComposeAggregateTestCheckFunc(tests...)
}
func testCheckHostConnectorItemValues(resourceName string, i int, connectorRequest api.CreateConnectorRequest) resource.TestCheckFunc {
	prefix := fmt.Sprintf("connector.%d.", i)
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, prefix+"id"),
		resource.TestCheckResourceAttr(resourceName, prefix+"name", connectorRequest.Name),
		resource.TestCheckResourceAttr(resourceName, prefix+"description", connectorRequest.Description),
		resource.TestCheckResourceAttr(resourceName, prefix+"vpn_region_id", connectorRequest.VpnRegionId),
		resource.TestCheckResourceAttrSet(resourceName, prefix+"ip_v4_address"),
		resource.TestCheckResourceAttrSet(resourceName, prefix+"ip_v6_address"),
		resource.TestCheckResourceAttrSet(resourceName, prefix+"profile"),
	)
}

func resourceHostOutputConfig(name string, host api.CreateHostRequest) string {

	hostResource := fmt.Sprintf(`
provider "openvpn" {}

resource "openvpn_host" "%s" {
	name = "%s"
	description = "%s"
	domain = "%s"
	internet_access = "%s"
`, name, host.Name, host.Description, host.Domain, host.InternetAccess)

	// 	Add a default connector if provided
	if len(host.Connectors) > 0 {
		connector := host.Connectors[0]
		hostResource += fmt.Sprintf(`
	connector {
		name = "%s"
		description = "%s"
		vpn_region_id = "%s"
	}`,
			connector.Name, connector.Description, connector.VpnRegionId)

	}

	hostResource += "\n}"

	return hostResource
}
