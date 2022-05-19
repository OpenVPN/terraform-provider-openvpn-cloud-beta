package api

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"
)

type NetworkItemType string

type ConnectionStatus string

type Connector struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	IpV4Address      string           `json:"ipV4Address"`
	IpV6Address      string           `json:"ipV6Address"`
	NetworkItemId    string           `json:"networkItemId"`
	NetworkItemType  NetworkItemType  `json:"networkItemType"`
	VpnRegionId      string           `json:"vpnRegionId"`
	ConnectionStatus ConnectionStatus `json:"connectionStatus"`
}

type CreateConnectorData struct {
	Name            string
	Description     string
	VpnRegionId     string
	NetworkItemId   string
	NetworkItemType NetworkItemType
}

type CreateConnectorRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	VpnRegionId string `json:"vpnRegionId"`
}

const (
	ConnectorByIdEndpoint        = "/connectors/%s"
	ConnectorsEndpoint           = "/connectors"
	ConnectorProfileByIdEndpoint = "/connectors/%s/profile"
)

const (
	NetworkItemTypeHost    NetworkItemType = "HOST"
	NetworkItemTypeNetwork NetworkItemType = "NETWORK"
)

const (
	ConnectionStatusOffline ConnectionStatus = "OFFLINE"
	ConnectionStatusOnline  ConnectionStatus = "ONLINE"
)

var NetworkItemPossibleValues = []string{string(NetworkItemTypeHost), string(NetworkItemTypeNetwork)}

func (c *Client) GetConnector(ctx context.Context, connectorId string) (*Connector, error) {
	connector := new(Connector)
	err := c.newRequest(ctx, "GET", c.apiEndpoint(ConnectorByIdEndpoint, connectorId), nil, connector)
	if err != nil {
		return nil, err
	}

	return connector, nil
}

func (c *Client) CreateConnector(ctx context.Context, request *CreateConnectorData) (*Connector, error) {

	createConnectorUrl := c.createConnectorUrl(c.apiEndpoint(ConnectorsEndpoint), request.NetworkItemId, request.NetworkItemType)

	connector := new(Connector)
	err := c.newRequestJSON(ctx, "POST", createConnectorUrl.String(), request.internalRequest(), connector)
	if err != nil {
		return nil, err
	}

	return connector, nil
}

func (c *Client) UpdateConnector(ctx context.Context, connectorID string, request *CreateConnectorData) (*Connector, error) {
	createConnectorUrl := c.createConnectorUrl(c.apiEndpoint(ConnectorByIdEndpoint, connectorID), request.NetworkItemId, request.NetworkItemType)

	connector := new(Connector)
	err := c.newRequestJSON(ctx, "PUT", createConnectorUrl.String(), request.internalRequest(), connector)
	if err != nil {
		return nil, err
	}

	return connector, nil
}

func (c *Client) DeleteConnector(ctx context.Context, networkItemId string, networkItemType NetworkItemType, connectorId string) error {
	endpoint := c.createConnectorUrl(c.apiEndpoint(ConnectorByIdEndpoint, connectorId), networkItemId, networkItemType)

	err := c.newRequest(ctx, "DELETE", endpoint.String(), nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetConnectorProfile(ctx context.Context, connectorId string) (string, error) {
	endpoint := c.apiEndpoint(ConnectorProfileByIdEndpoint, connectorId)
	response, err := c.newRequestWithResponse(ctx, "POST", endpoint, bytes.NewBufferString(""))
	if err != nil {
		return "", err
	}

	data, err := c.getBytesResponse(response)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *Client) createConnectorUrl(rawUrl, networkItemId string, networkItemType NetworkItemType) *url.URL {
	createConnectorUrl, _ := url.Parse(rawUrl)

	query := createConnectorUrl.Query()
	query.Set("networkItemId", networkItemId)
	query.Set("networkItemType", string(networkItemType))
	createConnectorUrl.RawQuery = query.Encode()

	return createConnectorUrl
}

func (t NetworkItemType) Validate() error {
	for _, possibleValue := range NetworkItemPossibleValues {
		if string(t) == possibleValue {
			return nil
		}
	}
	possibleValues := strings.Join(NetworkItemPossibleValues, ", ")
	return fmt.Errorf("invalid value for NetworkItemType: '%s'. Possible values are: %s", t, possibleValues)
}

func (r *CreateConnectorData) internalRequest() *CreateConnectorRequest {
	return &CreateConnectorRequest{
		Name:        r.Name,
		Description: r.Description,
		VpnRegionId: r.VpnRegionId,
	}
}
