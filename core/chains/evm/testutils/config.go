package testutils

import (
	"testing"

	"github.com/goplugin/plugin-common/pkg/logger"

	"github.com/goplugin/pluginv3.0/v2/core/chains/evm/config"
	"github.com/goplugin/pluginv3.0/v2/core/chains/evm/config/toml"
	"github.com/goplugin/pluginv3.0/v2/core/chains/evm/utils/big"
)

func NewTestChainScopedConfig(t testing.TB, overrideFn func(c *toml.EVMConfig)) config.ChainScopedConfig {
	var chainID = (*big.Big)(FixtureChainID)
	evmCfg := &toml.EVMConfig{
		ChainID: chainID,
		Chain:   toml.Defaults(chainID),
	}

	if overrideFn != nil {
		// We need to get the chainID from the override function first to load the correct chain defaults.
		// Then we apply the override values on top
		overrideFn(evmCfg)
		evmCfg.Chain = toml.Defaults(evmCfg.ChainID)
		overrideFn(evmCfg)
	}

	return config.NewTOMLChainScopedConfig(evmCfg, logger.Test(t))
}
