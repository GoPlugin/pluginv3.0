package ccipdata

import (
	cciptypes "github.com/goplugin/plugin-common/pkg/types/ccip"
)

const (
	ManuallyExecute = "manuallyExecute"
)

type OffRampReader interface {
	cciptypes.OffRampReader
}
