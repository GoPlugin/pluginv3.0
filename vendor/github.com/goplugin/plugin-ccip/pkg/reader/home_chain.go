package reader

import (
	"time"

	"github.com/goplugin/plugin-common/pkg/logger"
	"github.com/goplugin/plugin-common/pkg/types"

	reader_internal "github.com/goplugin/plugin-ccip/internal/reader"
)

type HomeChain = reader_internal.HomeChain

type ChainConfig = reader_internal.ChainConfig

type ChainConfigInfo = reader_internal.ChainConfigInfo

type OCR3ConfigWithMeta = reader_internal.OCR3ConfigWithMeta

type OCR3Config = reader_internal.OCR3Config

type OCR3Node = reader_internal.OCR3Node

func NewHomeChainReader(
	homeChainReader types.ContractReader,
	lggr logger.Logger,
	pollingInterval time.Duration,
	ccipConfigBoundContract types.BoundContract,
) HomeChain {
	return reader_internal.NewHomeChainConfigPoller(homeChainReader, lggr, pollingInterval, ccipConfigBoundContract)
}
