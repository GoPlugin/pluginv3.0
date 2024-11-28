package plugin

import (
	coscfg "github.com/goplugin/plugin-cosmos/pkg/cosmos/config"
	"github.com/goplugin/plugin-solana/pkg/solana"
	stkcfg "github.com/goplugin/plugin-starknet/relayer/pkg/plugin/config"

	"github.com/goplugin/pluginv3.0/v2/core/chains/evm/config/toml"
	"github.com/goplugin/pluginv3.0/v2/core/config"
)

//go:generate mockery --quiet --name GeneralConfig --output ./mocks/ --case=underscore

type GeneralConfig interface {
	config.AppConfig
	toml.HasEVMConfigs
	CosmosConfigs() coscfg.TOMLConfigs
	SolanaConfigs() solana.TOMLConfigs
	StarknetConfigs() stkcfg.TOMLConfigs
	// ConfigTOML returns both the user provided and effective configuration as TOML.
	ConfigTOML() (user, effective string)
}
