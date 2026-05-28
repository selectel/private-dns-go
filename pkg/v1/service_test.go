package v1

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var expectedServiceDetails = &ServiceDetails{
	Service: Service{
		ID:                "test_service",
		Project:           "test_project",
		NetworkID:         "test_network",
		HighAvailability:  true,
		IsRecursorEnabled: true,
	},
	Addresses: []*ServiceAddress{{Address: "192.168.0.1", CIDR: "192.168.0.0/24"}},
}

var expectedServicesList = []*Service{{
	ID:                "test_service",
	Project:           "test_project",
	NetworkID:         "test_network",
	HighAvailability:  true,
	IsRecursorEnabled: true,
}}

func ptr[T any](v T) *T {
	ptr := new(T)
	*ptr = v

	return ptr
}

func TestPrivateDNSClient__ListServices(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiServicesListJSON)),
				StatusCode: http.StatusOK,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)

		zones, err := client.ListServices(context.Background())
		require.NoError(t, err)

		require.Equal(t, expectedServicesList, zones)

		require.Equal(t, "http://test.com/services", httpClient.request.URL.String())
		require.Equal(t, http.MethodGet, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})

	t.Run("ApiErrorr", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiErrorJSON)),
				StatusCode: http.StatusNotFound,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)

		_, err := client.ListServices(context.Background())
		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "Not found", apiErr.Msg)
		require.Equal(t, apiErrorJSON, apiErr.Raw())

		require.Equal(t, "http://test.com/services", httpClient.request.URL.String())
		require.Equal(t, http.MethodGet, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})
}

func TestPrivateDNSClient__GetService(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiServiceDetailJSON)),
				StatusCode: http.StatusOK,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)
		service, err := client.GetService(context.Background(), "test_service")
		require.NoError(t, err)

		require.Equal(t, expectedServiceDetails, service)

		require.Equal(t, "http://test.com/services/test_service", httpClient.request.URL.String())
		require.Equal(t, http.MethodGet, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})

	t.Run("ApiErrorr", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiErrorJSON)),
				StatusCode: http.StatusNotFound,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)

		_, err := client.GetService(context.Background(), "test_service")
		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "Not found", apiErr.Msg)
		require.Equal(t, apiErrorJSON, apiErr.Raw())

		require.Equal(t, "http://test.com/services/test_service", httpClient.request.URL.String())
		require.Equal(t, http.MethodGet, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})
}

func TestPrivateDNSClient__DeleteService(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader("")),
				StatusCode: http.StatusNoContent,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)
		err := client.DeleteService(context.Background(), "test_service")
		require.NoError(t, err)

		require.Equal(t, "http://test.com/services/test_service", httpClient.request.URL.String())
		require.Equal(t, http.MethodDelete, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})

	t.Run("ApiErrorr", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiErrorJSON)),
				StatusCode: http.StatusNotFound,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)

		err := client.DeleteService(context.Background(), "test_service")
		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "Not found", apiErr.Msg)
		require.Equal(t, apiErrorJSON, apiErr.Raw())

		require.Equal(t, "http://test.com/services/test_service", httpClient.request.URL.String())
		require.Equal(t, http.MethodDelete, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})
}

func TestPrivateDNSClient__CreateService(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiServiceDetailJSON)),
				StatusCode: http.StatusCreated,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)

		dto := &ServiceCreateDTO{
			NetworkID:         "test_network",
			IsRecursorEnabled: ptr(true),
		}

		service, err := client.CreateService(context.Background(), dto)
		require.NoError(t, err)

		require.Equal(t, expectedServiceDetails, service)

		require.Equal(t, "http://test.com/services", httpClient.request.URL.String())
		require.Equal(t, http.MethodPost, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})

	t.Run("ApiErrorr", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiErrorJSON)),
				StatusCode: http.StatusNotFound,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)
		dto := &ServiceCreateDTO{
			NetworkID: "test_network",
		}

		_, err := client.CreateService(context.Background(), dto)
		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "Not found", apiErr.Msg)
		require.Equal(t, apiErrorJSON, apiErr.Raw())

		require.Equal(t, "http://test.com/services", httpClient.request.URL.String())
		require.Equal(t, http.MethodPost, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})
}

func TestPrivateDNSClient__UpdateService(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiServiceDetailJSON)),
				StatusCode: http.StatusNoContent,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)

		dto := &ServiceUpdateDTO{
			ServiceID:         "test_network",
			IsRecursorEnabled: ptr(true),
		}

		err := client.UpdateService(context.Background(), dto)
		require.NoError(t, err)

		require.Equal(t, "http://test.com/services/test_network", httpClient.request.URL.String())
		require.Equal(t, http.MethodPut, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})

	t.Run("ApiErrorr", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiErrorJSON)),
				StatusCode: http.StatusNotFound,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)
		dto := &ServiceUpdateDTO{
			ServiceID:         "test_network",
			IsRecursorEnabled: ptr(true),
		}

		err := client.UpdateService(context.Background(), dto)
		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "Not found", apiErr.Msg)
		require.Equal(t, apiErrorJSON, apiErr.Raw())

		require.Equal(t, "http://test.com/services/test_network", httpClient.request.URL.String())
		require.Equal(t, http.MethodPut, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})
}
