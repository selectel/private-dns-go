package v1

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var expectedZoneData = &ZoneDetails{
	Zone: Zone{
		ID:           "test_zone",
		Name:         "example.com.",
		Domain:       "example.com.",
		Project:      "test_project",
		TTL:          3600,
		SerialNumber: 0,
	},
	Records:  make([]*Record, 0),
	Bindings: make([]*ZoneBindings, 0),
}

var expectedZonesList = []*Zone{{
	ID:           "test_zone",
	Name:         "example.com.",
	Domain:       "example.com.",
	Project:      "test_project",
	TTL:          3600,
	SerialNumber: 9,
}}

func TestPrivateDNSClient__ListZones(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiZonesListJSON)),
				StatusCode: http.StatusOK,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)

		zonesQuery := &ZonesQuery{
			Domain: "some.domain",
		}

		zones, err := client.ListZones(context.Background(), zonesQuery)
		require.NoError(t, err)

		require.Equal(t, expectedZonesList, zones)

		require.Equal(t, "http://test.com/zones?domain=some.domain", httpClient.request.URL.String())
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

		_, err := client.ListZones(context.Background(), nil)
		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "Not found", apiErr.Msg)
		require.Equal(t, apiErrorJSON, apiErr.Raw())

		require.Equal(t, "http://test.com/zones", httpClient.request.URL.String())
		require.Equal(t, http.MethodGet, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})
}

func TestPrivateDNSClient__GetZone(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiZonesDetailsJSON)),
				StatusCode: http.StatusOK,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)
		zones, err := client.GetZone(context.Background(), "test_zone")
		require.NoError(t, err)

		require.Equal(t, expectedZoneData, zones)

		require.Equal(t, "http://test.com/zones/test_zone", httpClient.request.URL.String())
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

		_, err := client.GetZone(context.Background(), "test_zone")
		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "Not found", apiErr.Msg)
		require.Equal(t, apiErrorJSON, apiErr.Raw())

		require.Equal(t, "http://test.com/zones/test_zone", httpClient.request.URL.String())
		require.Equal(t, http.MethodGet, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})
}

func TestPrivateDNSClient__DeleteZone(t *testing.T) {
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
		err := client.DeleteZone(context.Background(), "test_zone")
		require.NoError(t, err)

		require.Equal(t, "http://test.com/zones/test_zone", httpClient.request.URL.String())
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

		err := client.DeleteZone(context.Background(), "test_zone")
		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "Not found", apiErr.Msg)
		require.Equal(t, apiErrorJSON, apiErr.Raw())

		require.Equal(t, "http://test.com/zones/test_zone", httpClient.request.URL.String())
		require.Equal(t, http.MethodDelete, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})
}

func TestPrivateDNSClient__CreateZone(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiZonesDetailsJSON)),
				StatusCode: http.StatusCreated,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)

		dto := &ZoneCreateDTO{
			Name:   "example.com.",
			Domain: "example.com,",
		}

		zones, err := client.CreateZone(context.Background(), dto)
		require.NoError(t, err)

		require.Equal(t, expectedZoneData, zones)

		require.Equal(t, "http://test.com/zones", httpClient.request.URL.String())
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
		dto := &ZoneCreateDTO{
			Name:   "example.com.",
			Domain: "example.com,",
		}
		_, err := client.CreateZone(context.Background(), dto)
		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "Not found", apiErr.Msg)
		require.Equal(t, apiErrorJSON, apiErr.Raw())

		require.Equal(t, "http://test.com/zones", httpClient.request.URL.String())
		require.Equal(t, http.MethodPost, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})
}

func TestPrivateDNSClient__UpdateZone(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader("")),
				StatusCode: http.StatusOK,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)

		err := client.UpdateZone(context.Background(), "test_zone", &ZoneUpdateDto{})
		require.NoError(t, err)

		require.Equal(t, "http://test.com/zones/test_zone", httpClient.request.URL.String())
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

		err := client.UpdateZone(context.Background(), "test_zone", &ZoneUpdateDto{})
		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "Not found", apiErr.Msg)
		require.Equal(t, apiErrorJSON, apiErr.Raw())

		require.Equal(t, "http://test.com/zones/test_zone", httpClient.request.URL.String())
		require.Equal(t, http.MethodPut, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})
}
