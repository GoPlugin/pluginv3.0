package node

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/goplugin/plugin-testing-framework/lib/blockchain"
	"github.com/goplugin/plugin-testing-framework/lib/utils/ptr"

	"github.com/goplugin/plugin-common/pkg/config"

	"github.com/goplugin/pluginv3.0/integration-tests/types/config/node"
	itutils "github.com/goplugin/pluginv3.0/integration-tests/utils"
	evmcfg "github.com/goplugin/pluginv3.0/v2/core/chains/evm/config/toml"
	ubig "github.com/goplugin/pluginv3.0/v2/core/chains/evm/utils/big"
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
)

func NewConfigFromToml(tomlConfig []byte, opts ...node.NodeConfigOpt) (*plugin.Config, error) {
	var cfg plugin.Config
	err := config.DecodeTOML(bytes.NewReader(tomlConfig), &cfg)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &cfg, nil
}

func WithPrivateEVMs(networks []blockchain.EVMNetwork, commonChainConfig *evmcfg.Chain, chainSpecificConfig map[int64]evmcfg.Chain) node.NodeConfigOpt {
	var evmConfigs []*evmcfg.EVMConfig
	for _, network := range networks {
		var evmNodes []*evmcfg.Node
		for i := range network.URLs {
			evmNodes = append(evmNodes, &evmcfg.Node{
				Name:    ptr.Ptr(fmt.Sprintf("%s-%d", network.Name, i)),
				WSURL:   itutils.MustURL(network.URLs[i]),
				HTTPURL: itutils.MustURL(network.HTTPURLs[i]),
			})
		}
		evmConfig := &evmcfg.EVMConfig{
			ChainID: ubig.New(big.NewInt(network.ChainID)),
			Nodes:   evmNodes,
			Chain:   evmcfg.Chain{},
		}
		if commonChainConfig != nil {
			evmConfig.Chain = *commonChainConfig
		}
		if chainSpecificConfig == nil {
			if overriddenChainCfg, ok := chainSpecificConfig[network.ChainID]; ok {
				evmConfig.Chain = overriddenChainCfg
			}
		}
		if evmConfig.Chain.FinalityDepth == nil && network.FinalityDepth > 0 {
			evmConfig.Chain.FinalityDepth = ptr.Ptr(uint32(network.FinalityDepth))
		}
		if evmConfig.Chain.FinalityTagEnabled == nil && network.FinalityTag {
			evmConfig.Chain.FinalityTagEnabled = ptr.Ptr(network.FinalityTag)
		}
		evmConfigs = append(evmConfigs, evmConfig)
	}
	return func(c *plugin.Config) {
		c.EVM = evmConfigs
	}
}
