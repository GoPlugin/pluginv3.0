package migration

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/goplugin/plugin-testing-framework/lib/logging"
	"github.com/goplugin/pluginv3.0/integration-tests/actions"
	"github.com/goplugin/pluginv3.0/integration-tests/docker/test_env"

	tc "github.com/goplugin/pluginv3.0/integration-tests/testconfig"
)

func TestVersionUpgrade(t *testing.T) {
	t.Parallel()

	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig([]string{"Migration"}, tc.Node)
	require.NoError(t, err, "Error getting config")

	err = config.PluginUpgradeImage.Validate()
	require.NoError(t, err, "Error validating upgrade image")

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestConfig(&config).
		WithTestInstance(t).
		WithStandardCleanup().
		WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
		WithCLNodes(1).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	// just restarting CL container with the same name, DB is still the same
	//
	// [Database]
	// MigrateOnStartup = true
	//
	// by default
	err = env.ClCluster.Nodes[0].UpgradeVersion(*config.PluginUpgradeImage.Image, *config.PluginUpgradeImage.Version)
	require.NoError(t, err)

}
