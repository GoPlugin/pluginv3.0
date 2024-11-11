package soak

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/goplugin/plugin-testing-framework/logging"

	"github.com/goplugin/pluginv3.0/integration-tests/actions"
	tc "github.com/goplugin/pluginv3.0/integration-tests/testconfig"
	"github.com/goplugin/pluginv3.0/integration-tests/testsetups"
)

func TestForwarderOCRSoak(t *testing.T) {
	l := logging.GetTestLogger(t)
	// Use this variable to pass in any custom EVM specific TOML values to your Plugin nodes
	customNetworkTOML := `[EVM.Transactions]
ForwardersEnabled = true`
	// Uncomment below for debugging TOML issues on the node
	// fmt.Println("Using Plugin TOML\n---------------------")
	// fmt.Println(networks.AddNetworkDetailedConfig(config.BaseOCRP2PV1Config, customNetworkTOML, network))
	// fmt.Println("---------------------")

	config, err := tc.GetConfig("Soak", tc.OCR)
	require.NoError(t, err, "Error getting config")

	ocrSoakTest, err := testsetups.NewOCRSoakTest(t, &config, true)
	require.NoError(t, err, "Error creating soak test")
	ocrSoakTest.DeployEnvironment(customNetworkTOML, &config)
	if ocrSoakTest.Environment().WillUseRemoteRunner() {
		return
	}
	t.Cleanup(func() {
		if err := actions.TeardownRemoteSuite(ocrSoakTest.TearDownVals(t)); err != nil {
			l.Error().Err(err).Msg("Error tearing down environment")
		}
	})
	ocrSoakTest.Setup(&config)
	ocrSoakTest.Run()
}
