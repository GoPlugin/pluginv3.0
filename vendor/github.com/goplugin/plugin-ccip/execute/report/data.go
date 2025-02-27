package report

import (
	"sort"

	"github.com/goplugin/plugin-ccip/execute/exectypes"
	cciptypes "github.com/goplugin/plugin-ccip/pkg/types/ccipocr3"
)

// markNewMessagesExecuted compares an execute plugin report with the commit report metadata and marks the new messages
// as executed.
func markNewMessagesExecuted(
	execReport cciptypes.ExecutePluginReportSingleChain, report exectypes.CommitData,
) exectypes.CommitData {
	// Mark new messages executed.
	for i := 0; i < len(execReport.Messages); i++ {
		report.ExecutedMessages =
			append(report.ExecutedMessages, execReport.Messages[i].Header.SequenceNumber)
	}
	sort.Slice(
		report.ExecutedMessages,
		func(i, j int) bool { return report.ExecutedMessages[i] < report.ExecutedMessages[j] })

	return report
}
