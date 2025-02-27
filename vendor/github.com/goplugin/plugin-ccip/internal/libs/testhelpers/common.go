package testhelpers

import (
	"github.com/goplugin/plugin-libocr/commontypes"
	libocrtypes "github.com/goplugin/plugin-libocr/ragep2p/types"

	rmntypes "github.com/goplugin/plugin-ccip/commit/merkleroot/rmn/types"
	"github.com/goplugin/plugin-ccip/internal/libs/testhelpers/rand"
	"github.com/goplugin/plugin-ccip/internal/reader"
	"github.com/goplugin/plugin-ccip/pkg/types/ccipocr3"
)

func SetupConfigInfo(chainSelector ccipocr3.ChainSelector,
	readers []libocrtypes.PeerID,
	fChain uint8,
	cfg []byte) reader.ChainConfigInfo {
	return reader.ChainConfigInfo{
		ChainSelector: chainSelector,
		ChainConfig: reader.HomeChainConfigMapper{
			Readers: readers,
			FChain:  fChain,
			Config:  cfg,
		},
	}
}

func CreateOracleIDToP2pID(ids ...int) map[commontypes.OracleID]libocrtypes.PeerID {
	res := make(map[commontypes.OracleID]libocrtypes.PeerID)
	for _, id := range ids {
		res[commontypes.OracleID(id)] = libocrtypes.PeerID{byte(id)}
	}
	return res
}

func CreateRMNRemoteCfg() rmntypes.RemoteConfig {
	return rmntypes.RemoteConfig{
		ContractAddress: rand.RandomBytes(20),
		ConfigDigest:    rand.RandomBytes32(),
		Signers: []rmntypes.RemoteSignerInfo{
			{
				OnchainPublicKey: rand.RandomBytes(20),
				NodeIndex:        rand.RandomUint64(),
			},
		},
		MinSigners:       rand.RandomUint64(),
		ConfigVersion:    rand.RandomUint32(),
		RmnReportVersion: rand.RandomReportVersion(),
	}
}
