package protocol

import (
	"github.com/goplugin/plugin-libocr/commontypes"
	"github.com/goplugin/plugin-libocr/offchainreporting2plus/types"
)

type TelemetrySender interface {
	RoundStarted(
		configDigest types.ConfigDigest,
		epoch uint32,
		round uint8,
		leader commontypes.OracleID,
	)
}
