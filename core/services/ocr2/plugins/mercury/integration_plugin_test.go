//go:build integration

package mercury_test

import (
	"testing"

	"github.com/goplugin/pluginv3.0/v2/core/config/env"
)

func TestIntegration_MercuryV1_Plugin(t *testing.T) {
	t.Setenv(string(env.MercuryPlugin.Cmd), "plugin-mercury")
	integration_MercuryV1(t)
}

func TestIntegration_MercuryV2_Plugin(t *testing.T) {
	t.Setenv(string(env.MercuryPlugin.Cmd), "plugin-mercury")
	integration_MercuryV2(t)
}

func TestIntegration_MercuryV3_Plugin(t *testing.T) {
	t.Setenv(string(env.MercuryPlugin.Cmd), "plugin-mercury")
	integration_MercuryV3(t)
}
