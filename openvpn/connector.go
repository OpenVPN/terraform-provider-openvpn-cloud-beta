package openvpn

import (
	"context"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-openvpn/openvpn/api"
)

func resourceConnector() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectorCreate,
		ReadContext:   resourceConnectorRead,
		UpdateContext: resourceConnectorUpdate,
		DeleteContext: resourceConnectorDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
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
				Required: true,
			},
			"network_item_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: func(i interface{}, path cty.Path) diag.Diagnostics {
					value := api.NetworkItemType(i.(string))
					return diag.FromErr(value.Validate())
				},
			},
			"vpn_region_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"profile": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceConnectorUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, ok := i.(*api.Client)
	if !ok {
		return diag.Errorf("invalid api client")
	}

	connectorId, ok := data.Get("id").(string)
	if !ok {
		return diag.Errorf("invalid connector id")
	}

	request := &api.CreateConnectorData{
		Name:            data.Get("name").(string),
		Description:     data.Get("description").(string),
		VpnRegionId:     data.Get("vpn_region_id").(string),
		NetworkItemId:   data.Get("network_item_id").(string),
		NetworkItemType: api.NetworkItemType(data.Get("network_item_type").(string)),
	}

	connector, err := client.UpdateConnector(ctx, connectorId, request)
	if err != nil {
		return diag.FromErr(err)
	}

	diagnostics := setConnectorData(data, connector)
	if diagnostics != nil {
		return diagnostics
	}

	connectorProfile, err := client.GetConnectorProfile(ctx, connector.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("profile", connectorProfile)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceConnectorDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, ok := i.(*api.Client)
	if !ok {
		return diag.Errorf("invalid api client")
	}

	networkItemId, ok := data.Get("network_item_id").(string)
	if !ok {
		return diag.Errorf("invalid network_item_id")
	}

	networkItemType, ok := data.Get("network_item_type").(string)
	if !ok {
		return diag.Errorf("invalid network_item_type")
	}

	connectorId, ok := data.Get("id").(string)
	if !ok {
		return diag.Errorf("invalid connector id")
	}

	err := client.DeleteConnector(ctx, networkItemId, api.NetworkItemType(networkItemType), connectorId)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceConnectorCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, ok := i.(*api.Client)
	if !ok {
		return diag.Errorf("invalid api client")
	}

	request := &api.CreateConnectorData{
		Name:            data.Get("name").(string),
		Description:     data.Get("description").(string),
		VpnRegionId:     data.Get("vpn_region_id").(string),
		NetworkItemId:   data.Get("network_item_id").(string),
		NetworkItemType: api.NetworkItemType(data.Get("network_item_type").(string)),
	}

	connector, err := client.CreateConnector(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	diagnostics := setConnectorData(data, connector)
	if diagnostics != nil {
		return diagnostics
	}

	connectorProfile, err := client.GetConnectorProfile(ctx, connector.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("profile", connectorProfile)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceConnectorRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
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

	diagnostics := setConnectorData(data, connector)
	if diagnostics != nil {
		return diagnostics
	}

	connectorProfile, err := client.GetConnectorProfile(ctx, connector.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("profile", connectorProfile)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setConnectorData(data *schema.ResourceData, connector *api.Connector) diag.Diagnostics {
	data.SetId(connector.ID)
	err := data.Set("name", connector.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("description", connector.Description)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("ip_v4_address", connector.IpV4Address)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("ip_v6_address", connector.IpV6Address)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("network_item_id", connector.NetworkItemId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("network_item_type", connector.NetworkItemType)
	if err != nil {
		return diag.FromErr(err)
	}
	err = data.Set("vpn_region_id", connector.VpnRegionId)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
