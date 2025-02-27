package core

import (
	"context"

	"github.com/goplugin/plugin-libocr/offchainreporting2/reportingplugin/median"

	"github.com/goplugin/plugin-common/pkg/types"
)

type PluginMedian interface {
	// NewMedianFactory returns a new ReportingPluginFactory. If provider implements GRPCClientConn, it can be forwarded efficiently via proxy.
	NewMedianFactory(ctx context.Context, provider types.MedianProvider, contractID string, dataSource, juelsPerFeeCoin, gasPriceSubunits median.DataSource, errorLog ErrorLog) (types.ReportingPluginFactory, error)
}
