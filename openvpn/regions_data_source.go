package openvpn

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-openvpn/openvpn/api"
)

func dataSourceRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataRegionsRead,
		Schema: map[string]*schema.Schema{
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"continent": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"country_iso": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataRegionsRead(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, ok := i.(*api.Client)
	if !ok {
		return diag.Errorf("invalid api client")
	}
	regions, err := client.ListRegions(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	regionsData := make([]map[string]interface{}, len(regions))

	combinedID := ""

	for i, region := range regions {
		regionsData[i] = map[string]interface{}{
			"id":          region.ID,
			"continent":   region.Continent,
			"country":     region.Country,
			"country_iso": region.CountryISO,
			"region_name": region.RegionName,
		}

		combinedID += region.ID
	}

	err = d.Set("regions", regionsData)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(combinedID)

	return nil
}
