package smoke

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/goplugin/pluginv3.0/integration-tests/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/goplugin/plugin-testing-framework/lib/logging"
	"github.com/goplugin/plugin-testing-framework/lib/utils/testcontext"

	"github.com/goplugin/pluginv3.0/integration-tests/actions"
	"github.com/goplugin/pluginv3.0/integration-tests/contracts"
	"github.com/goplugin/pluginv3.0/integration-tests/docker/test_env"
	tc "github.com/goplugin/pluginv3.0/integration-tests/testconfig"
)

func TestForwarderOCRBasic(t *testing.T) {
	t.Parallel()
	l := logging.GetTestLogger(t)

	config, err := tc.GetConfig([]string{"Smoke"}, tc.ForwarderOcr)
	require.NoError(t, err, "Error getting config")

	privateNetwork, err := actions.EthereumNetworkConfigFromConfig(l, &config)
	require.NoError(t, err, "Error building ethereum network config")

	env, err := test_env.NewCLTestEnvBuilder().
		WithTestInstance(t).
		WithTestConfig(&config).
		WithPrivateEthereumNetwork(privateNetwork.EthereumNetworkConfig).
		WithMockAdapter().
		WithCLNodes(6).
		WithStandardCleanup().
		Build()
	require.NoError(t, err)

	nodeClients := env.ClCluster.NodeAPIs()
	bootstrapNode, workerNodes := nodeClients[0], nodeClients[1:]

	workerNodeAddresses, err := actions.PluginNodeAddressesLocal(workerNodes)
	require.NoError(t, err, "Retreiving on-chain wallet addresses for plugin nodes shouldn't fail")

	evmNetwork, err := env.GetFirstEvmNetwork()
	require.NoError(t, err, "Error getting first evm network")

	sethClient, err := utils.TestAwareSethClient(t, config, evmNetwork)
	require.NoError(t, err, "Error getting seth client")

	err = actions.FundPluginNodesFromRootAddress(l, sethClient, contracts.PluginClientToPluginNodeWithKeysAndAddress(env.ClCluster.NodeAPIs()), big.NewFloat(*config.Common.PluginNodeFunding))
	require.NoError(t, err, "Failed to fund the nodes")

	t.Cleanup(func() {
		// ignore error, we will see failures in the logs anyway
		_ = actions.ReturnFundsFromNodes(l, sethClient, contracts.PluginClientToPluginNodeWithKeysAndAddress(env.ClCluster.NodeAPIs()))
	})

	linkContract, err := actions.LinkTokenContract(l, sethClient, config.OCR)
	require.NoError(t, err, "Error loading/deploying link token contract")

	fundingAmount := big.NewFloat(.05)
	l.Info().Str("ETH amount per node", fundingAmount.String()).Msg("Funding Plugin nodes")
	err = actions.FundPluginNodesFromRootAddress(l, sethClient, contracts.PluginClientToPluginNodeWithKeysAndAddress(workerNodes), fundingAmount)
	require.NoError(t, err, "Error funding Plugin nodes")

	operators, authorizedForwarders, _ := actions.DeployForwarderContracts(
		t, sethClient, common.HexToAddress(linkContract.Address()), len(workerNodes),
	)

	require.Equal(t, len(workerNodes), len(operators), "Number of operators should match number of worker nodes")

	for i := range workerNodes {
		actions.AcceptAuthorizedReceiversOperator(
			t, l, sethClient, operators[i], authorizedForwarders[i], []common.Address{workerNodeAddresses[i]},
		)
		require.NoError(t, err, "Accepting Authorize Receivers on Operator shouldn't fail")
		actions.TrackForwarder(t, sethClient, authorizedForwarders[i], workerNodes[i])
	}

	ocrInstances, err := actions.DeployOCRContractsForwarderFlow(
		l,
		sethClient,
		config.OCR,
		common.HexToAddress(linkContract.Address()),
		contracts.PluginClientToPluginNodeWithKeysAndAddress(workerNodes),
		authorizedForwarders,
	)
	require.NoError(t, err, "Error deploying OCR contracts")

	err = actions.CreateOCRJobsWithForwarderLocal(ocrInstances, bootstrapNode, workerNodes, 5, env.MockAdapter, fmt.Sprint(sethClient.ChainID))
	require.NoError(t, err, "failed to setup forwarder jobs")
	err = actions.WatchNewOCRRound(l, sethClient, 1, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(10*time.Minute))
	require.NoError(t, err, "error watching for new OCR round")

	answer, err := ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Getting latest answer from OCR contract shouldn't fail")
	require.Equal(t, int64(5), answer.Int64(), "Expected latest answer from OCR contract to be 5 but got %d", answer.Int64())

	err = actions.SetAllAdapterResponsesToTheSameValueLocal(10, ocrInstances, workerNodes, env.MockAdapter)
	require.NoError(t, err)
	err = actions.WatchNewOCRRound(l, sethClient, 2, contracts.V1OffChainAgrregatorToOffChainAggregatorWithRounds(ocrInstances), time.Duration(10*time.Minute))
	require.NoError(t, err, "error watching for new OCR round")

	answer, err = ocrInstances[0].GetLatestAnswer(testcontext.Get(t))
	require.NoError(t, err, "Error getting latest OCR answer")
	require.Equal(t, int64(10), answer.Int64(), "Expected latest answer from OCR contract to be 10 but got %d", answer.Int64())
}
