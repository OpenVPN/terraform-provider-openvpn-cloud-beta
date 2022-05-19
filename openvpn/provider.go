package openvpn

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"net/http"
	"terraform-provider-openvpn/openvpn/api"
)

const ProviderName = "openvpn"

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("OVPN_HOST", nil),
				ValidateFunc: validation.IsURLWithHTTPS,
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVPN_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OVPN_CLIENT_SECRET", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"openvpn_host":      resourceHost(),
			"openvpn_connector": resourceConnector(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"openvpn_regions":   dataSourceRegions(),
			"openvpn_host":      dataSourceHost(),
			"openvpn_connector": dataSourceConnector(),
		},
		ConfigureContextFunc: configureProviderContext,
	}
}

func configureProviderContext(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	authConfig := &api.AuthConfig{
		Host:         data.Get("host").(string),
		ClientID:     data.Get("client_id").(string),
		ClientSecret: data.Get("client_secret").(string),
	}

	httpClient := &http.Client{}

	client := api.NewClient(httpClient, authConfig)
	err := client.Authenticate(ctx)
	if err != nil {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "authentication failed: " + err.Error(),
				Detail:   err.Error(),
			},
		}
	}

	return client, nil
}
