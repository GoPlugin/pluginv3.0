package web_test

import (
	"net/http"
	"testing"

	"github.com/goplugin/pluginv3.0/v2/core/internal/cltest"
	"github.com/goplugin/pluginv3.0/v2/core/internal/testutils/configtest"
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
)

func TestCors_DefaultOrigins(t *testing.T) {
	t.Parallel()

	config := configtest.NewGeneralConfig(t, func(c *plugin.Config, s *plugin.Secrets) {
		c.WebServer.AllowOrigins = ptr("http://localhost:3000,http://localhost:6689")
	})

	tests := []struct {
		origin     string
		statusCode int
	}{
		{"http://localhost:3000", http.StatusOK},
		{"http://localhost:6689", http.StatusOK},
		{"http://localhost:1234", http.StatusForbidden},
	}

	for _, test := range tests {
		t.Run(test.origin, func(t *testing.T) {
			app := cltest.NewApplicationWithConfig(t, config)

			client := app.NewHTTPClient(nil)

			headers := map[string]string{"Origin": test.origin}
			resp, cleanup := client.Get("/v2/chains/evm", headers)
			defer cleanup()
			cltest.AssertServerResponse(t, resp, test.statusCode)
		})
	}
}

func TestCors_OverrideOrigins(t *testing.T) {
	t.Parallel()

	tests := []struct {
		allow      string
		origin     string
		statusCode int
	}{
		{"http://plugin.com", "http://plugin.com", http.StatusOK},
		{"http://plugin.com", "http://localhost:3000", http.StatusForbidden},
		{"*", "http://plugin.com", http.StatusOK},
		{"*", "http://localhost:3000", http.StatusOK},
	}

	for _, test := range tests {
		t.Run(test.origin, func(t *testing.T) {
			config := configtest.NewGeneralConfig(t, func(c *plugin.Config, s *plugin.Secrets) {
				c.WebServer.AllowOrigins = ptr(test.allow)
			})
			app := cltest.NewApplicationWithConfig(t, config)

			client := app.NewHTTPClient(nil)

			headers := map[string]string{"Origin": test.origin}
			resp, cleanup := client.Get("/v2/chains/evm", headers)
			defer cleanup()
			cltest.AssertServerResponse(t, resp, test.statusCode)
		})
	}
}
