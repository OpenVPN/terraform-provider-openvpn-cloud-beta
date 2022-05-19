package openvpn

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var defaultProviderFactory = map[string]func() (*schema.Provider, error){
	ProviderName: defaultProvider,
}

var defaultProviderInstance = Provider()

func defaultProvider() (*schema.Provider, error) {
	return defaultProviderInstance, nil
}

func TestProvider(t *testing.T) {
	provider := Provider()
	assert.NotNil(t, provider)
}

func TestProvider_Validate(t *testing.T) {
	provider := Provider()
	err := provider.InternalValidate()
	assert.NoError(t, err)
}

func TestConfigureProvider(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Only for acceptance testing")
	}
	testAccPreCheck(t)

	ctx := context.Background()

	provider := Provider()

	data := map[string]interface{}{
		"host":          os.Getenv("OVPN_HOST"),
		"client_id":     os.Getenv("OVPN_CLIENT_ID"),
		"client_secret": os.Getenv("OVPN_CLIENT_SECRET"),
	}
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, data)

	info, diag := configureProviderContext(ctx, resourceData)
	assert.NotNil(t, info)
	if !assert.False(t, diag.HasError()){
		t.Log(diag)
	}
}

func testAccPreCheck(t *testing.T) {
	validateEnvVar(t, "OVPN_HOST")
	validateEnvVar(t, "OVPN_CLIENT_ID")
	validateEnvVar(t, "OVPN_CLIENT_SECRET")
}

func validateEnvVar(t *testing.T, envVar string) {
	require.NotEmptyf(t, os.Getenv(envVar), "%s must be set for acceptance tests", envVar)
}
