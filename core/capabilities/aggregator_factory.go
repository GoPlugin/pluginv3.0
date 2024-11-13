package capabilities

import (
	"fmt"

	"github.com/goplugin/plugin-common/pkg/capabilities/consensus/ocr3/datafeeds"
	"github.com/goplugin/plugin-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/goplugin/plugin-common/pkg/logger"
	"github.com/goplugin/plugin-common/pkg/values"
	"github.com/goplugin/pluginv3.0/v2/core/capabilities/streams"
)

func NewAggregator(name string, config values.Map, lggr logger.Logger) (types.Aggregator, error) {
	switch name {
	case "data_feeds":
		mc := streams.NewCodec(lggr)
		return datafeeds.NewDataFeedsAggregator(config, mc, lggr)
	default:
		return nil, fmt.Errorf("aggregator %s not supported", name)
	}
}
