package v1

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrivateDNSClient__processApiError(t *testing.T) {
	t.Run("Correct", func(t *testing.T) {
		client := &PrivateDNSClient{}
		response := &http.Response{
			Body:       io.NopCloser(strings.NewReader(apiErrorJSON)),
			StatusCode: http.StatusInternalServerError,
		}

		err := client.processAPIError(response)

		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "Not found", apiErr.Msg)
		require.Equal(t, apiErrorJSON, apiErr.Raw())
		require.Equal(t, response.StatusCode, apiErr.Code)
	})

	t.Run("IncorrectJSON", func(t *testing.T) {
		client := &PrivateDNSClient{}
		data := `{"test":some}`
		response := &http.Response{
			Body:       io.NopCloser(strings.NewReader(data)),
			StatusCode: http.StatusInternalServerError,
		}

		err := client.processAPIError(response)

		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "", apiErr.Msg)
		require.Equal(t, data, apiErr.Raw())
		require.Equal(t, response.StatusCode, apiErr.Code)
	})

	t.Run("Plain", func(t *testing.T) {
		client := &PrivateDNSClient{}
		data := `Plain`
		response := &http.Response{
			Body:       io.NopCloser(strings.NewReader(data)),
			StatusCode: http.StatusInternalServerError,
		}

		err := client.processAPIError(response)

		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "", apiErr.Msg)
		require.Equal(t, data, apiErr.Raw())
		require.Equal(t, response.StatusCode, apiErr.Code)
	})
}

func TestPrivateDNSClient__doRequest(t *testing.T) {
	t.Run("WithTarget", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(`{"test":"some"}`)),
				StatusCode: http.StatusOK,
			},
		}
		cfg := &Config{
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://test.com", nil)
		require.NoError(t, err)

		target := map[string]string{}
		err = client.doRequest(req, 200, &target)
		require.NoError(t, err)

		require.Equal(t, "some", target["test"])
	})

	t.Run("WithTransportError", func(t *testing.T) {
		httpClient := &testHTTPClient{
			err: errors.New("kaboom"),
		}
		cfg := &Config{
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://test.com", nil)
		require.NoError(t, err)

		target := map[string]string{}
		err = client.doRequest(req, 200, &target)

		var tErr *TransportErr
		require.ErrorAs(t, err, &tErr)
		require.ErrorIs(t, err, httpClient.err)
	})

	t.Run("WithApiError", func(t *testing.T) {
		httpClient := &testHTTPClient{
			response: &http.Response{
				Body:       io.NopCloser(strings.NewReader(apiErrorJSON)),
				StatusCode: http.StatusInternalServerError,
			},
		}
		cfg := &Config{
			HTTPClient: httpClient,
		}

		client := NewPrivateDNSClient(cfg)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://test.com", nil)
		require.NoError(t, err)

		target := map[string]string{}
		err = client.doRequest(req, 200, &target)

		var apiErr *APIErr
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "Not found", apiErr.Msg)
		require.Equal(t, apiErrorJSON, apiErr.Raw())
		require.Equal(t, httpClient.response.StatusCode, apiErr.Code)
	})
}

func TestPrivateDNSClient__makeRequest(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cfg := &Config{
			AuthToken: "test_token",
			URL:       "http://test.com",
		}
		client := NewPrivateDNSClient(cfg)

		testQuery := url.Values{}
		testQuery.Add("foo", "bar")
		req, err := client.makeRequest(context.Background(), http.MethodGet, "/some", map[string]string{"test": "test"}, testQuery)
		require.NoError(t, err)

		require.Equal(t, cfg.AuthToken, req.Header.Get(xAuthHeader))
		require.Equal(t, moduleUserAgent, req.Header.Get(userAgentHeader))
		require.Equal(t, req.URL.String(), "http://test.com/some?foo=bar")
	})
	t.Run("withUserAgent", func(t *testing.T) {
		cfg := &Config{
			AuthToken: "test_token",
			URL:       "http://test.com",
			UserAgent: "some-custom-user-agent",
		}
		client := NewPrivateDNSClient(cfg)

		testQuery := url.Values{}
		testQuery.Add("foo", "bar")
		req, err := client.makeRequest(context.Background(), http.MethodGet, "/some", map[string]string{"test": "test"}, testQuery)
		require.NoError(t, err)

		require.Equal(t, cfg.AuthToken, req.Header.Get(xAuthHeader))
		require.Equal(t, cfg.UserAgent+" "+moduleUserAgent, req.Header.Get(userAgentHeader))
		require.Equal(t, req.URL.String(), "http://test.com/some?foo=bar")
	})
}
