package fluxmonitorv2

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/goplugin/pluginv3.0/v2/core/null"
)

// FluxMonitorRoundStatsV2 defines the stats for a round
type FluxMonitorRoundStatsV2 struct {
	ID              uint64
	PipelineRunID   null.Int64
	Aggregator      common.Address
	RoundID         uint32
	NumNewRoundLogs uint64
	NumSubmissions  uint64
}
