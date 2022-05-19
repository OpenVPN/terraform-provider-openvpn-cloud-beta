package api

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockHttpClient struct {
	mock.Mock
}

func newMockHttpClient() *mockHttpClient {
	return &mockHttpClient{}
}

func (m *mockHttpClient) Do(request *http.Request) (*http.Response, error) {
	args := m.Called(request)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *mockHttpClient) mockDo(t *testing.T, responseBody interface{}, requestTestFunction func(request *http.Request)) *mock.Call {
	response := createTestResponse(t, responseBody)

	call := m.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(response.Result(), nil)

	if requestTestFunction != nil {
		call = call.Run(func(args mock.Arguments) {
			request := getRequestFromArgs(t, args)
			requestTestFunction(request)
		})
	}

	return call
}

func (m *mockHttpClient) mockDoBytes(t *testing.T, responseBody []byte, requestTestFunction func(request *http.Request)) *mock.Call {
	response := httptest.NewRecorder()
	response.Body = bytes.NewBuffer(responseBody)

	call := m.
		On("Do", mock.AnythingOfType("*http.Request")).
		Return(response.Result(), nil)

	if requestTestFunction != nil {
		call = call.Run(func(args mock.Arguments) {
			request := getRequestFromArgs(t, args)
			requestTestFunction(request)
		})
	}

	return call
}

func createTestResponse(t *testing.T, mockResponseBody interface{}) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	if mockResponseBody != nil {
		data, err := json.Marshal(mockResponseBody)
		require.NoError(t, err)
		response.Body = bytes.NewBuffer(data)
	}
	return response
}
