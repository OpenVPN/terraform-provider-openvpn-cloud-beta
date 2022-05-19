package openvpn

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-openvpn/openvpn/api"
)

func resourceHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHostCreate,
		ReadContext:   resourceHostRead,
		UpdateContext: resourceHostUpdate,
		DeleteContext: resourceHostDelete,
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
			"internet_access": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"system_subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"connector": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
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
						"vpn_region_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"profile": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

// CREATE
func resourceHostCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, ok := i.(*api.Client)
	if !ok {
		return diag.Errorf("invalid api client")
	}

	request := makeHostCreationRequest(data)

	host, err := client.CreateHost(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(host.ID)
	err = data.Set("system_subnets", host.SystemSubnets)
	if err != nil {
		return diag.FromErr(err)
	}

	diagnostics := setConnectorsList(ctx, data, client, host.Connectors)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func makeHostCreationRequest(data *schema.ResourceData) *api.CreateHostRequest {
	request := &api.CreateHostRequest{
		Name:           data.Get("name").(string),
		Description:    data.Get("description").(string),
		Domain:         data.Get("domain").(string),
		InternetAccess: data.Get("internet_access").(string),
	}

	connectorsI := data.Get("connector").([]interface{})
	request.Connectors = make([]api.CreateConnectorRequest, len(connectorsI))
	for i, connectorI := range connectorsI {
		connectorData := connectorI.(map[string]interface{})
		request.Connectors[i] = api.CreateConnectorRequest{
			Name:        connectorData["name"].(string),
			Description: connectorData["description"].(string),
			VpnRegionId: connectorData["vpn_region_id"].(string),
		}
	}
	return request
}

// READ
func resourceHostRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
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

	diagnostics := setConnectorsList(ctx, data, client, []api.Connector{host.Connectors[0]})
	if diagnostics != nil {
		return diagnostics
	}
	return nil
}

// UPDATE
func resourceHostUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, ok := i.(*api.Client)
	if !ok {
		return diag.Errorf("invalid api client")
	}

	hostID, ok := data.Get("id").(string)
	if !ok {
		return diag.Errorf("invalid id")
	}

	request := &api.UpdateHostRequest{
		Name:           data.Get("name").(string),
		Description:    data.Get("description").(string),
		Domain:         data.Get("domain").(string),
		InternetAccess: data.Get("internet_access").(string),
	}

	host, err := client.UpdateHost(ctx, hostID, request)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(host.ID)
	err = data.Set("system_subnets", host.SystemSubnets)
	if err != nil {
		return diag.FromErr(err)
	}

	connectors, diagnostics := updateConnectors(ctx, data, client, host.ID)
	if diagnostics != nil {
		return diagnostics
	}

	diagnostics = setConnectorsList(ctx, data, client, connectors)
	if diagnostics != nil {
		return diagnostics
	}

	return nil
}

func updateConnectors(ctx context.Context, data *schema.ResourceData, client *api.Client, hostID string) ([]api.Connector, diag.Diagnostics) {
	connectorI := data.Get("connector").([]interface{})[0]
	connectorData := connectorI.(map[string]interface{})
	connectorRequest := &api.CreateConnectorData{
		Name:            connectorData["name"].(string),
		Description:     connectorData["description"].(string),
		VpnRegionId:     connectorData["vpn_region_id"].(string),
		NetworkItemId:   hostID,
		NetworkItemType: api.NetworkItemTypeHost,
	}

	connector, err := client.UpdateConnector(ctx, connectorData["id"].(string), connectorRequest)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	connectors := []api.Connector{*connector}
	return connectors, nil
}

// DELETE
func resourceHostDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, ok := i.(*api.Client)
	if !ok {
		return diag.Errorf("invalid api client")
	}

	hostID, ok := data.Get("id").(string)
	if !ok {
		return diag.Errorf("invalid id")
	}

	err := client.DeleteHost(ctx, hostID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setConnectorsList(ctx context.Context, data *schema.ResourceData, client *api.Client, connectors []api.Connector) diag.Diagnostics {
	connectorsList := make([]interface{}, len(connectors))
	for i, connector := range connectors {
		connectorsData, err := getConnectorsListItem(ctx, client, connector)
		if err != nil {
			return diag.FromErr(err)
		}
		connectorsList[i] = connectorsData
	}
	err := data.Set("connector", connectorsList)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getConnectorsListItem(ctx context.Context, client *api.Client, connector api.Connector) (map[string]interface{}, error) {
	connectorsData := map[string]interface{}{
		"id":            connector.ID,
		"name":          connector.Name,
		"description":   connector.Description,
		"vpn_region_id": connector.VpnRegionId,
		"ip_v4_address": connector.IpV4Address,
		"ip_v6_address": connector.IpV6Address,
	}

	connectorProfile, err := client.GetConnectorProfile(ctx, connector.ID)
	if err != nil {
		return nil, err
	}
	connectorsData["profile"] = connectorProfile
	return connectorsData, nil
}
