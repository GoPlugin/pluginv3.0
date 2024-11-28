package capabilities

import (
	"fmt"

	"github.com/goplugin/plugin-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/goplugin/plugin-common/pkg/logger"
	"github.com/goplugin/plugin-common/pkg/values"
	"github.com/goplugin/pluginv3.0/v2/core/services/relay/evm"
)

func NewEncoder(name string, config *values.Map, lggr logger.Logger) (types.Encoder, error) {
	switch name {
	case "EVM":
		return evm.NewEVMEncoder(config)
	// TODO: add a "no-op" encoder for users who only want to use dynamic ones?
	// https://smartcontract-it.atlassian.net/browse/CAPPL-88
	default:
		return nil, fmt.Errorf("encoder %s not supported", name)
	}
}
