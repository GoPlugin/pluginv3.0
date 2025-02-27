package llo

import (
	"github.com/goplugin/plugin-libocr/offchainreporting2plus/ocr3types"

	llotypes "github.com/goplugin/plugin-common/pkg/types/llo"
)

type ChannelDefinitionWithID struct {
	llotypes.ChannelDefinition
	ChannelID llotypes.ChannelID
}

type ChannelHash [32]byte

type Transmitter interface {
	// NOTE: Mercury doesn't actually transmit on-chain, so there is no
	// "contract" involved with the transmitter.
	// - Transmit should be implemented and send to Mercury server
	// - FromAccount() should return CSA public key
	ocr3types.ContractTransmitter[llotypes.ReportInfo]
}
