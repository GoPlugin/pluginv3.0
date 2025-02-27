package tokenprice

import (
	"context"
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"

	"github.com/goplugin/plugin-common/pkg/logger"

	"github.com/goplugin/plugin-libocr/commontypes"
	"github.com/goplugin/plugin-libocr/offchainreporting2plus/types"

	"github.com/goplugin/plugin-ccip/internal/plugincommon"
	"github.com/goplugin/plugin-ccip/internal/reader"
	pkgreader "github.com/goplugin/plugin-ccip/pkg/reader"
	cciptypes "github.com/goplugin/plugin-ccip/pkg/types/ccipocr3"
	"github.com/goplugin/plugin-ccip/pluginconfig"
)

type processor struct {
	oracleID         commontypes.OracleID
	lggr             logger.Logger
	offChainCfg      pluginconfig.CommitOffchainConfig
	destChain        cciptypes.ChainSelector
	chainSupport     plugincommon.ChainSupport
	tokenPriceReader pkgreader.PriceReader
	homeChain        reader.HomeChain
	fRoleDON         int
}

func NewProcessor(
	oracleID commontypes.OracleID,
	lggr logger.Logger,
	offChainCfg pluginconfig.CommitOffchainConfig,
	destChain cciptypes.ChainSelector,
	chainSupport plugincommon.ChainSupport,
	tokenPriceReader pkgreader.PriceReader,
	homeChain reader.HomeChain,
	fRoleDON int,
) plugincommon.PluginProcessor[Query, Observation, Outcome] {
	return &processor{
		oracleID:         oracleID,
		lggr:             lggr,
		offChainCfg:      offChainCfg,
		destChain:        destChain,
		chainSupport:     chainSupport,
		tokenPriceReader: tokenPriceReader,
		homeChain:        homeChain,
		fRoleDON:         fRoleDON,
	}
}

func (p *processor) Query(ctx context.Context, prevOutcome Outcome) (Query, error) {
	return Query{}, nil
}

func (p *processor) ValidateObservation(
	prevOutcome Outcome,
	query Query,
	ao plugincommon.AttributedObservation[Observation],
) error {
	return validateObservedTokenPrices(ao.Observation.FeedTokenPrices)
}

func (p *processor) Outcome(
	ctx context.Context,
	_ Outcome,
	_ Query,
	aos []plugincommon.AttributedObservation[Observation],
) (Outcome, error) {
	p.lggr.Infow("processing token price outcome")
	// If set to zero, no prices will be reported (i.e keystone feeds would be active).
	if p.offChainCfg.TokenPriceBatchWriteFrequency.Duration() == 0 {
		p.lggr.Debugw("TokenPriceBatchWriteFrequency is set to zero, no prices will be reported")
		return Outcome{}, nil
	}

	consensusObservation, err := p.getConsensusObservation(aos)
	if err != nil {
		return Outcome{}, err
	}

	tokenPriceOutcome := p.selectTokensForUpdate(consensusObservation)
	p.lggr.Infow(
		"outcome token prices",
		"token prices", tokenPriceOutcome,
	)
	return Outcome{
		TokenPrices: tokenPriceOutcome,
	}, nil
}

func (p *processor) Close() error {
	return nil
}

func validateObservedTokenPrices(tokenPrices []cciptypes.TokenPrice) error {
	tokensWithPrice := mapset.NewSet[types.Account]()
	for _, t := range tokenPrices {
		if tokensWithPrice.Contains(t.TokenID) {
			return fmt.Errorf("duplicate token price for token: %s", t.TokenID)
		}
		tokensWithPrice.Add(t.TokenID)

		if t.Price.IsEmpty() {
			return fmt.Errorf("token price must not be empty")
		}
	}
	return nil
}

var _ plugincommon.PluginProcessor[Query, Observation, Outcome] = &processor{}
