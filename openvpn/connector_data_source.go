package openvpn

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-openvpn/openvpn/api"
	"time"
)

func dataSourceConnector() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectorRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(2 * time.Minute),
		},
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
			"ip_v4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_v6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_item_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_item_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpn_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceConnectorRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, ok := i.(*api.Client)
	if !ok {
		return diag.Errorf("invalid api client")
	}

	connectorId, ok := data.Get("id").(string)
	if !ok {
		return diag.Errorf("invalid id")
	}

	connector, err := client.GetConnector(ctx, connectorId)
	if err != nil {
		return diag.FromErr(err)
	}

	setConnectorData(data, connector)

	connectorProfile, err := client.GetConnectorProfile(ctx, connectorId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("profile", connectorProfile)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
