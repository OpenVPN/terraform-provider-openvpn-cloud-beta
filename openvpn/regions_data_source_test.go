package openvpn

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

const dataRegionsOutputConfig = `
provider "openvpn" {
}

data "openvpn_regions" "test" {
}
`

func TestDataRegions(t *testing.T) {
	dataSourceName := "data.openvpn_regions.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: defaultProviderFactory,
		Steps: []resource.TestStep{
			{
				Config: dataRegionsOutputConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "%", "2"),
					resource.TestCheckResourceAttrSet(dataSourceName, "regions.#"),
					resource.TestCheckResourceAttr(dataSourceName, "regions.0.%", "5"),
					testCheckResourceAttributesMapSet(dataSourceName, "regions.0",
						"id", "continent", "country", "country_iso", "region_name"),
				),
			},
		},
	})
}

func testCheckResourceAttributesMapSet(name, attrName string, keys ...string) resource.TestCheckFunc {
	numKeys := len(keys)
	testCheckFuncs := make([]resource.TestCheckFunc, numKeys)
	for i, key := range keys {
		testCheckFuncs[i] = resource.TestCheckResourceAttrSet(name, fmt.Sprintf("%s.%s", attrName, key))
	}
	return resource.ComposeAggregateTestCheckFunc(testCheckFuncs...)
}

