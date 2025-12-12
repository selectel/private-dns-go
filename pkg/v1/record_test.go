package v1

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var expectedRecordList = []*Record{{
	Type:      "A",
	Domain:    "sub.example.com.",
	TTL:       -1,
	Values:    []string{"192.168.0.1"},
	Generated: true,
}}

func TestPrivateDNSClient__GetRecords(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiRecordListJSON)),
				StatusCode: http.StatusOK,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)

		query := &RecordsQuery{
			Domain: "some.domain",
		}

		result, err := client.ListRecords(context.Background(), "test_zone", query)
		require.NoError(t, err)

		require.Equal(t, expectedRecordList, result)

		require.Equal(t, "http://test.com/zones/test_zone/recordset?domain=some.domain", httpClient.request.URL.String())
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

		_, err := client.ListRecords(context.Background(), "test_zone", nil)
		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "Not found", apiErr.Msg)
		require.Equal(t, apiErrorJSON, apiErr.Raw())

		require.Equal(t, "http://test.com/zones/test_zone/recordset", httpClient.request.URL.String())
		require.Equal(t, http.MethodGet, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})
}

func TestPrivateDNSClient__PutRecords(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiRecordListJSON)),
				StatusCode: http.StatusOK,
			},
		}
		cfg := &Config{
			AuthToken:  "testToken",
			URL:        "http://test.com",
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)

		result, err := client.PutRecords(context.Background(), "test_zone", &PutRecordsDTO{})
		require.NoError(t, err)

		require.Equal(t, expectedRecordList, result)

		require.Equal(t, "http://test.com/zones/test_zone/recordset", httpClient.request.URL.String())
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

		_, err := client.PutRecords(context.Background(), "test_zone", &PutRecordsDTO{})
		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "Not found", apiErr.Msg)
		require.Equal(t, apiErrorJSON, apiErr.Raw())

		require.Equal(t, "http://test.com/zones/test_zone/recordset", httpClient.request.URL.String())
		require.Equal(t, http.MethodPut, httpClient.request.Method)
		require.Equal(t, cfg.AuthToken, httpClient.request.Header.Get(xAuthHeader))
	})
}
