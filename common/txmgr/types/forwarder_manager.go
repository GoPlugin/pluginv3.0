package types

import (
	"github.com/goplugin/plugin-common/pkg/services"
	"github.com/goplugin/pluginv3.0/v2/common/types"
)

//go:generate mockery --quiet --name ForwarderManager --output ./mocks/ --case=underscore
type ForwarderManager[ADDR types.Hashable] interface {
	services.Service
	ForwarderFor(addr ADDR) (forwarder ADDR, err error)
	// Converts payload to be forwarder-friendly
	ConvertPayload(dest ADDR, origPayload []byte) ([]byte, error)
}
