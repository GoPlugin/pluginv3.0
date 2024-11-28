package client_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/goplugin/plugin-common/pkg/logger"

	"github.com/goplugin/pluginv3.0/v2/core/chains/evm/client"
	"github.com/goplugin/pluginv3.0/v2/core/chains/evm/config/chaintype"
	"github.com/goplugin/pluginv3.0/v2/core/chains/evm/testutils"
)

func TestNewEvmClient(t *testing.T) {
	t.Parallel()

	noNewHeadsThreshold := 3 * time.Minute
	selectionMode := ptr("HighestHead")
	leaseDuration := 0 * time.Second
	pollFailureThreshold := ptr(uint32(5))
	pollInterval := 10 * time.Second
	syncThreshold := ptr(uint32(5))
	nodeIsSyncingEnabled := ptr(false)
	chainTypeStr := ""
	finalizedBlockOffset := ptr[uint32](16)
	enforceRepeatableRead := ptr(true)
	deathDeclarationDelay := time.Second * 3
	noNewFinalizedBlocksThreshold := time.Second * 5
	finalizedBlockPollInterval := time.Second * 4
	newHeadsPollInterval := time.Second * 4
	nodeConfigs := []client.NodeConfig{
		{
			Name:    ptr("foo"),
			WSURL:   ptr("ws://foo.test"),
			HTTPURL: ptr("http://foo.test"),
		},
	}
	finalityDepth := ptr(uint32(10))
	finalityTagEnabled := ptr(true)
	chainCfg, nodePool, nodes, err := client.NewClientConfigs(selectionMode, leaseDuration, chainTypeStr, nodeConfigs,
		pollFailureThreshold, pollInterval, syncThreshold, nodeIsSyncingEnabled, noNewHeadsThreshold, finalityDepth,
		finalityTagEnabled, finalizedBlockOffset, enforceRepeatableRead, deathDeclarationDelay, noNewFinalizedBlocksThreshold,
		finalizedBlockPollInterval, newHeadsPollInterval)
	require.NoError(t, err)

	client := client.NewEvmClient(nodePool, chainCfg, nil, logger.Test(t), testutils.FixtureChainID, nodes, chaintype.ChainType(chainTypeStr))
	require.NotNil(t, client)
}