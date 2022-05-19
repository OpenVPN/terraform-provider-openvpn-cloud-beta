package openvpn

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-openvpn/openvpn/api"
)

func dataSourceHost() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHostRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_access": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"connectors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpn_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceHostRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, ok := i.(*api.Client)
	if !ok {
		return diag.Errorf("invalid api client")
	}
	hostID, ok := data.Get("id").(string)
	if !ok {
		return diag.Errorf("invalid id")
	}
	host, err := client.GetHost(ctx, hostID)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(host.ID)
	err = data.Set("name", host.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("description", host.Description)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("internet_access", host.InternetAccess)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("domain", host.Domain)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("system_subnets", host.SystemSubnets)
	if err != nil {
		return diag.FromErr(err)
	}

	connectorsData := make([]map[string]interface{}, len(host.Connectors))
	for i, connector := range host.Connectors {
		connectorsData[i] = map[string]interface{}{
			"id": connector.ID,
			"name": connector.Name,
			"vpn_region_id": connector.VpnRegionId,
		}
	}
	err = data.Set("connectors", connectorsData)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
