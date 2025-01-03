package solidity_cross_tests_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/goplugin/pluginv3.0/v2/core/services/vrf/solidity_cross_tests"
)

var (
	keyHash   = secretKey.PublicKey.MustHash()
	jobID     = common.BytesToHash([]byte("1234567890abcdef1234567890abcdef"))
	seed      = big.NewInt(1)
	sender    = common.HexToAddress("0xecfcab0a285d3380e488a39b4bb21e777f8a4eac")
	fee       = big.NewInt(100)
	requestID = common.HexToHash("0xcafe")
	raw       = solidity_cross_tests.RawRandomnessRequestLog{
		KeyHash:   keyHash,
		Seed:      seed,
		JobID:     jobID,
		Sender:    sender,
		Fee:       fee,
		RequestID: requestID,
		Raw: types.Log{
			// A raw, on-the-wire RandomnessRequestLog is the concat of fields as uint256's
			Data: append(append(append(append(
				keyHash.Bytes(),
				common.BigToHash(seed).Bytes()...),
				common.BytesToHash(sender.Bytes()).Bytes()...),
				common.BigToHash(fee).Bytes()...),
				requestID.Bytes()...),
			Topics: []common.Hash{{}, jobID},
		},
	}
)

func TestVRFParseRandomnessRequestLog(t *testing.T) {
	r := solidity_cross_tests.RawRandomnessRequestLogToRandomnessRequestLog(&raw)
	rawLog, err := r.RawData()
	require.NoError(t, err)
	assert.Equal(t, rawLog, raw.Raw.Data)
	nR, err := solidity_cross_tests.ParseRandomnessRequestLog(types.Log{
		Data:   rawLog,
		Topics: []common.Hash{solidity_cross_tests.VRFRandomnessRequestLogTopic(), jobID},
	})
	require.NoError(t, err)
	require.True(t, r.Equal(*nR),
		"Round-tripping RandomnessRequestLog through serialization and parsing "+
			"resulted in a different log.")
}
