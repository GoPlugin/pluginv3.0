package protocol

import (
	"github.com/goplugin/plugin-libocr/commontypes"
	"github.com/goplugin/plugin-libocr/offchainreporting2plus/types"
)

type TelemetrySender interface {
	RoundStarted(
		configDigest types.ConfigDigest,
		epoch uint64,
		seqNr uint64,
		round uint64,
		leader commontypes.OracleID,
	)
}
