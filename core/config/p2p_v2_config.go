package config

import (
	commonconfig "github.com/goplugin/plugin-common/pkg/config"

	ocrcommontypes "github.com/goplugin/plugin-libocr/commontypes"
)

type V2 interface {
	Enabled() bool
	AnnounceAddresses() []string
	DefaultBootstrappers() (locators []ocrcommontypes.BootstrapperLocator)
	DeltaDial() commonconfig.Duration
	DeltaReconcile() commonconfig.Duration
	ListenAddresses() []string
}
