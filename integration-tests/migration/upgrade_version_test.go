package migration

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/goplugin/pluginv3.0/integration-tests/docker/test_env"

	tc "github.com/goplugin/pluginv3.0/integration-tests/testconfig"
)

func TestVersionUpgrade(t *testing.T) {
	t.Parallel()

	config, err := tc.GetConfig("Migration", tc.Node)
	require.NoError(t, err, "Error getting config")

	err = config.PluginUpgradeImage.Validate()
	require.NoError(t, err, "Error validating upgrade image")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestConfig(&config).
		WithTestInstance(t).
		WithStandardCleanup().
		WithGeth().
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
