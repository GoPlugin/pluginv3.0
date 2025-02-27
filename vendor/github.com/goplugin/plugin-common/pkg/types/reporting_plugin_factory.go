package types

import (
	libocr "github.com/goplugin/plugin-libocr/offchainreporting2plus/types"
)

type ReportingPluginFactory interface {
	Service
	libocr.ReportingPluginFactory
}
