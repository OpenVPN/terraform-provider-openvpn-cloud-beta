package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	client     HttpClient
	authData   *AuthData
	authConfig *AuthConfig
}

type HttpClient interface {
	Do(request *http.Request) (*http.Response, error)
}

type ErrorResponse struct {
	Errors      map[string][]string `json:"errors"`
	Path        string              `json:"path"`
	RequestId   string              `json:"requestId"`
	Status      int                 `json:"status"`
	StatusError string              `json:"statusError"`
	Timestamp   int64               `json:"timestamp"`
}

func NewClient(client HttpClient, authConfig *AuthConfig) *Client {
	return &Client{client: client, authConfig: authConfig}
}

func processJsonResponse(response *http.Response, body interface{}) error {
	err := processResponseError(response)
	if err != nil {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(body)
	if err != nil {
		return err
	}

	err = response.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getBytesResponse(response *http.Response) ([]byte, error) {
	err := processResponseError(response)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = response.Body.Close()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func processResponseError(response *http.Response) error {
	if response.StatusCode >= 400 {
		errorBody := &ErrorResponse{}
		err := json.NewDecoder(response.Body).Decode(errorBody)
		if err != nil {
			return fmt.Errorf("%s %s %s", response.Request.Method, response.Request.URL.Path, response.Status)
		}
		return errorBody
	}

	return nil
}

func (c Client) apiEndpoint(format string, a ...interface{}) string {
	return c.authConfig.apiUrl(format, a...)
}

func (c Client) IsAuthenticated() bool {
	return c.authData != nil
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%d %s %v", e.Status, e.StatusError, e.Errors)
}

func (c *Client) newRequestJSON(ctx context.Context, method, url string, reqBody, resBody interface{}) error {
	var reqBodyReader io.Reader
	if reqBody != nil {
		data, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}
		reqBodyReader = bytes.NewReader(data)
	}
	return c.newRequest(ctx, method, url, reqBodyReader, resBody)
}

func (c *Client) newRequest(ctx context.Context, method, url string, reqBodyReader io.Reader, resBody interface{}) error {
	response, err := c.newRequestWithResponse(ctx, method, url, reqBodyReader)
	if err != nil {
		return err
	}

	if resBody != nil {
		err = processJsonResponse(response, resBody)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) newRequestWithResponse(ctx context.Context, method string, url string, reqBodyReader io.Reader) (*http.Response, error) {
	if !c.IsAuthenticated() {
		return nil, errors.New("authentication is required")
	}

	request, err := http.NewRequest(method, url, reqBodyReader)
	if err != nil {
		return nil, err
	}
	c.authData.AuthorizeRequest(request)

	response, err := c.client.Do(request.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return response, nil
}
