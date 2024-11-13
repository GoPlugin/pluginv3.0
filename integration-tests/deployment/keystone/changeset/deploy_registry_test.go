package changeset_test

import (
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/require"

	"github.com/goplugin/plugin-common/pkg/logger"
	"github.com/goplugin/pluginv3.0/integration-tests/deployment"
	"github.com/goplugin/pluginv3.0/integration-tests/deployment/keystone/changeset"
	"github.com/goplugin/pluginv3.0/integration-tests/deployment/memory"
)

func TestDeployCapabilityRegistry(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	ab := deployment.NewMemoryAddressBook()
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	registrySel := env.AllChainSelectors()[0]
	resp, err := changeset.DeployCapabilityRegistry(lggr, env, ab, registrySel)
	require.NoError(t, err)
	require.NotNil(t, resp)
	// capabilities registry should be deployed on chain 0
	addrs, err := resp.AddressBook.AddressesForChain(registrySel)
	require.NoError(t, err)
	require.Len(t, addrs, 1)

	// no capabilities registry on chain 1
	require.NotEqual(t, registrySel, env.AllChainSelectors()[1])
	oaddrs, _ := resp.AddressBook.AddressesForChain(env.AllChainSelectors()[1])
	require.Len(t, oaddrs, 0)

}
