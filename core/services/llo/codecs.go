package llo

import (
	llotypes "github.com/goplugin/plugin-common/pkg/types/llo"
	"github.com/goplugin/plugin-data-streams/llo"

	"github.com/goplugin/pluginv3.0/v2/core/services/llo/evm"
)

// NOTE: All supported codecs must be specified here
func NewCodecs() map[llotypes.ReportFormat]llo.ReportCodec {
	codecs := make(map[llotypes.ReportFormat]llo.ReportCodec)

	codecs[llotypes.ReportFormatJSON] = llo.JSONReportCodec{}
	codecs[llotypes.ReportFormatEVMPremiumLegacy] = evm.ReportCodecPremiumLegacy{}

	return codecs
}
