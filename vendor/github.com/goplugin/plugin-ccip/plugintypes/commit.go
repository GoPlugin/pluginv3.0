package plugintypes

import (
	"time"

	cciptypes "github.com/goplugin/plugin-ccip/pkg/types/ccipocr3"
)

// NOTE: The following type should be moved to internal plugin types after it's not required anymore in plugin repo.
// Right now it's only used in a plugin repo test: TestCCIPReader_CommitReportsGTETimestamp

type CommitPluginReportWithMeta struct {
	Report    cciptypes.CommitPluginReport `json:"report"`
	Timestamp time.Time                    `json:"timestamp"`
	BlockNum  uint64                       `json:"blockNum"`
}
