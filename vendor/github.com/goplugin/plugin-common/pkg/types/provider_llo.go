package types

import (
	ocrtypes "github.com/goplugin/plugin-libocr/offchainreporting2plus/types"

	"github.com/goplugin/plugin-common/pkg/types/llo"
)

type LLOConfigProvider interface {
	OffchainConfigDigester() ocrtypes.OffchainConfigDigester
	// One instance will be run per config tracker
	ContractConfigTrackers() []ocrtypes.ContractConfigTracker
}

type LLOProvider interface {
	Service
	LLOConfigProvider
	ShouldRetireCache() llo.ShouldRetireCache
	ContractTransmitter() llo.Transmitter
	ChannelDefinitionCache() llo.ChannelDefinitionCache
}
