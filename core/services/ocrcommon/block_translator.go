package ocrcommon

import (
	"context"
	"math/big"

	"github.com/goplugin/plugin-common/pkg/logger"

	evmclient "github.com/goplugin/pluginv3.0/v2/core/chains/evm/client"
	"github.com/goplugin/pluginv3.0/v2/core/chains/evm/config/chaintype"
	evmtypes "github.com/goplugin/pluginv3.0/v2/core/chains/evm/types"
)

// BlockTranslator converts emitted block numbers (from block.number) into a
// block number range suitable for query in FilterLogs
type BlockTranslator interface {
	NumberToQueryRange(ctx context.Context, changedInL1Block uint64) (fromBlock *big.Int, toBlock *big.Int)
}

// NewBlockTranslator returns the block translator for the given chain
func NewBlockTranslator(cfg Config, client evmclient.Client, lggr logger.Logger) BlockTranslator {
	switch cfg.ChainType() {
	case chaintype.ChainArbitrum:
		return NewArbitrumBlockTranslator(client, lggr)
	case "", chaintype.ChainCelo, chaintype.ChainGnosis, chaintype.ChainKroma, chaintype.ChainMetis, chaintype.ChainOptimismBedrock, chaintype.ChainScroll, chaintype.ChainWeMix, chaintype.ChainXLayer, chaintype.ChainZkEvm, chaintype.ChainZkSync, chaintype.ChainZircuit:
		fallthrough
	default:
		return &l1BlockTranslator{}
	}
}

type l1BlockTranslator struct{}

func (*l1BlockTranslator) NumberToQueryRange(_ context.Context, changedInL1Block uint64) (fromBlock *big.Int, toBlock *big.Int) {
	return big.NewInt(int64(changedInL1Block)), big.NewInt(int64(changedInL1Block))
}

func (*l1BlockTranslator) OnNewLongestChain(context.Context, *evmtypes.Head) {}
