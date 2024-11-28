package changeset

import (
	"testing"

	chain_selectors "github.com/goplugin/chain-selectors"
	"go.uber.org/zap/zapcore"

	"github.com/goplugin/pluginv3.0/integration-tests/deployment"
	"github.com/goplugin/pluginv3.0/integration-tests/deployment/memory"
	"github.com/goplugin/pluginv3.0/v2/core/logger"
)

var (
	// TestChain is the chain used by the in-memory environment.
	TestChain = chain_selectors.Chain{
		EvmChainID: 90000001,
		Selector:   909606746561742123,
		Name:       "Test Chain",
		VarName:    "",
	}
)

func newMemoryEnv(t *testing.T) deployment.Environment {
	lggr := logger.TestLogger(t)
	memEnvConf := memory.MemoryEnvironmentConfig{
		Chains:         1,
		Nodes:          4,
		Bootstraps:     1,
		RegistryConfig: deployment.CapabilityRegistryConfig{},
	}
	return memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memEnvConf)
}
