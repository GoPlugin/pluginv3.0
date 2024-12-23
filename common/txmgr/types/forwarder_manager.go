package types

import (
	"context"

	"github.com/goplugin/plugin-common/pkg/services"

	"github.com/goplugin/pluginv3.0/v2/common/types"
)

type ForwarderManager[ADDR types.Hashable] interface {
	services.Service
	ForwarderFor(ctx context.Context, addr ADDR) (forwarder ADDR, err error)
	ForwarderForOCR2Feeds(ctx context.Context, eoa, ocr2Aggregator ADDR) (forwarder ADDR, err error)
	// Converts payload to be forwarder-friendly
	ConvertPayload(dest ADDR, origPayload []byte) ([]byte, error)
}
