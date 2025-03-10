package median

import (
	"context"
	"fmt"
	"math/big"

	"github.com/goplugin/plugin-libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/goplugin/plugin-libocr/offchainreporting2plus/types"

	"github.com/goplugin/plugin-common/pkg/types"
)

const typeName = "MedianReport"

type reportCodec struct {
	codec types.Codec
}

var _ median.ReportCodec = &reportCodec{}

func (r *reportCodec) BuildReport(ctx context.Context, observations []median.ParsedAttributedObservation) (ocrtypes.Report, error) {
	if len(observations) == 0 {
		return nil, fmt.Errorf("cannot build report from empty attributed observations")
	}

	return r.codec.Encode(ctx, aggregate(observations), typeName)
}

func (r *reportCodec) MedianFromReport(ctx context.Context, report ocrtypes.Report) (*big.Int, error) {
	agg := &aggregatedAttributedObservation{}
	if err := r.codec.Decode(ctx, report, agg, typeName); err != nil {
		return nil, err
	}
	observations := make([]*big.Int, len(agg.Observations))
	copy(observations, agg.Observations)
	medianObservation := len(agg.Observations) / 2
	return agg.Observations[medianObservation], nil
}

func (r *reportCodec) MaxReportLength(ctx context.Context, n int) (int, error) {
	return r.codec.GetMaxDecodingSize(ctx, n, typeName)
}
