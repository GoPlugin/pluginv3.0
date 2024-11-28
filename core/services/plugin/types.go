package plugin

import (
	coscfg "github.com/goplugin/plugin-cosmos/pkg/cosmos/config"
	solcfg "github.com/goplugin/plugin-solana/pkg/solana/config"
	stkcfg "github.com/goplugin/plugin-starknet/relayer/pkg/plugin/config"

	"github.com/goplugin/pluginv3.0/v2/core/chains/evm/config/toml"
	"github.com/goplugin/pluginv3.0/v2/core/config"
)

type GeneralConfig interface {
	config.AppConfig
	toml.HasEVMConfigs
	CosmosConfigs() coscfg.TOMLConfigs
	SolanaConfigs() solcfg.TOMLConfigs
	StarknetConfigs() stkcfg.TOMLConfigs
	AptosConfigs() RawConfigs
	// ConfigTOML returns both the user provided and effective configuration as TOML.
	ConfigTOML() (user, effective string)
}
