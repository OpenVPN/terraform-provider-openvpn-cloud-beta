package api

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
)

type Host struct {
	ID             string      `json:"id,omitempty"`
	Name           string      `json:"name,omitempty"`
	Description    string      `json:"description,omitempty"`
	InternetAccess string      `json:"internetAccess,omitempty"`
	Domain         string      `json:"domain,omitempty"`
	Connectors     []Connector `json:"connectors,omitempty"`
	SystemSubnets  []string    `json:"systemSubnets,omitempty"`
}

type CreateHostRequest struct {
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	Domain         string                   `json:"domain"`
	InternetAccess string                   `json:"internetAccess"`
	Connectors     []CreateConnectorRequest `json:"connectors"`
}

type UpdateHostRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Domain         string `json:"domain"`
	InternetAccess string `json:"internetAccess"`
}

const HostsEndpoint = "/hosts"
const HostsDetailsEndpoint = "/hosts/%s"

func (c Client) GetHost(ctx context.Context, id string) (*Host, error) {
	host := new(Host)

	err := c.newRequest(ctx, "GET", c.apiEndpoint(HostsDetailsEndpoint, id), nil, host)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func (c Client) CreateHost(ctx context.Context, createHostRequest *CreateHostRequest) (*Host, error) {
	host := new(Host)

	err := c.newRequestJSON(ctx, "POST", c.apiEndpoint(HostsEndpoint), createHostRequest, host)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func (c Client) UpdateHost(ctx context.Context, id string, updateHostRequest *UpdateHostRequest) (*Host, error) {
	host := new(Host)

	err := c.newRequestJSON(ctx, "PUT", c.apiEndpoint(HostsDetailsEndpoint, id), updateHostRequest, host)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func (c Client) DeleteHost(ctx context.Context, id string) error {
	return c.newRequest(ctx, "DELETE", c.apiEndpoint(HostsDetailsEndpoint, id), nil, nil)
}

func (c Client) AuthorizeRequest(request *http.Request) error {
	if !c.IsAuthenticated() {
		return errors.New("authentication is required")
	}
	c.authData.AuthorizeRequest(request)
	return nil
}
