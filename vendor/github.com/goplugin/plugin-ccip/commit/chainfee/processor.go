package chainfee

import (
	"context"

	"github.com/goplugin/plugin-ccip/internal/reader"
	readerpkg "github.com/goplugin/plugin-ccip/pkg/reader"
	"github.com/goplugin/plugin-ccip/pluginconfig"

	"github.com/goplugin/plugin-common/pkg/logger"

	"github.com/goplugin/plugin-ccip/internal/plugincommon"
	cciptypes "github.com/goplugin/plugin-ccip/pkg/types/ccipocr3"
)

type processor struct {
	destChain    cciptypes.ChainSelector
	lggr         logger.Logger
	homeChain    reader.HomeChain
	ccipReader   readerpkg.CCIPReader
	cfg          pluginconfig.CommitOffchainConfig
	chainSupport plugincommon.ChainSupport
	fRoleDON     int
}

func NewProcessor(
	lggr logger.Logger,
	destChain cciptypes.ChainSelector,
	homeChain reader.HomeChain,
	ccipReader readerpkg.CCIPReader,
	offChainConfig pluginconfig.CommitOffchainConfig,
	chainSupport plugincommon.ChainSupport,
	fRoleDON int,
) plugincommon.PluginProcessor[Query, Observation, Outcome] {
	return &processor{
		lggr:         lggr,
		destChain:    destChain,
		homeChain:    homeChain,
		ccipReader:   ccipReader,
		fRoleDON:     fRoleDON,
		chainSupport: chainSupport,
		cfg:          offChainConfig,
	}
}

func (p *processor) Query(ctx context.Context, prevOutcome Outcome) (Query, error) {
	return Query{}, nil
}

var _ plugincommon.PluginProcessor[Query, Observation, Outcome] = &processor{}

func (p *processor) Close() error {
	return nil
}
