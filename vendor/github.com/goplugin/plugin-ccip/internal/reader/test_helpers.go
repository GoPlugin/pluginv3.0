package reader

import (
	"github.com/goplugin/plugin-libocr/commontypes"
	libocrtypes "github.com/goplugin/plugin-libocr/ragep2p/types"

	cciptypes "github.com/goplugin/plugin-ccip/pkg/types/ccipocr3"
)

var (
	chainA       = cciptypes.ChainSelector(1)
	chainB       = cciptypes.ChainSelector(2)
	chainC       = cciptypes.ChainSelector(3)
	oracleAId    = commontypes.OracleID(1)
	p2pOracleAId = libocrtypes.PeerID{byte(oracleAId)}
	oracleBId    = commontypes.OracleID(2)
	p2pOracleBId = libocrtypes.PeerID{byte(oracleBId)}
	oracleCId    = commontypes.OracleID(3)
	p2pOracleCId = libocrtypes.PeerID{byte(oracleCId)}
)
