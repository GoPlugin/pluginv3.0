package smoke

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/goplugin/plugin-testing-framework/blockchain"
	ctf_config "github.com/goplugin/plugin-testing-framework/config"
	"github.com/goplugin/plugin-testing-framework/k8s/environment"
	"github.com/goplugin/plugin-testing-framework/k8s/pkg/helm/plugin"
	eth "github.com/goplugin/plugin-testing-framework/k8s/pkg/helm/ethereum"
	"github.com/goplugin/plugin-testing-framework/logging"
	"github.com/goplugin/plugin-testing-framework/networks"
	"github.com/goplugin/plugin-testing-framework/utils/testcontext"

	"github.com/goplugin/pluginv3.0/integration-tests/actions"
	"github.com/goplugin/pluginv3.0/integration-tests/actions/ocr2vrf_actions"
	"github.com/goplugin/pluginv3.0/integration-tests/actions/ocr2vrf_actions/ocr2vrf_constants"
	"github.com/goplugin/pluginv3.0/integration-tests/client"
	"github.com/goplugin/pluginv3.0/integration-tests/config"
	"github.com/goplugin/pluginv3.0/integration-tests/contracts"
	"github.com/goplugin/pluginv3.0/integration-tests/testconfig"
	tc "github.com/goplugin/pluginv3.0/integration-tests/testconfig"
)

var ocr2vrfSmokeConfig *testconfig.TestConfig

func TestOCR2VRFRedeemModel(t *testing.T) {
	t.Parallel()
	t.Skip("VRFv3 is on pause, skipping")
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig("Smoke", tc.OCR2)
	if err != nil {
		t.Fatal(err)
	}

	testEnvironment, testNetwork := setupOCR2VRFEnvironment(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment, l)
	require.NoError(t, err, "Error connecting to blockchain")
	contractDeployer, err := contracts.NewContractDeployer(chainClient, l)
	require.NoError(t, err, "Error building contract deployer")
	pluginNodes, err := client.ConnectPluginNodes(testEnvironment)
	require.NoError(t, err, "Error connecting to Plugin nodes")
	nodeAddresses, err := actions.PluginNodeAddresses(pluginNodes)
	require.NoError(t, err, "Retreiving on-chain wallet addresses for plugin nodes shouldn't fail")

	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, pluginNodes, nil, zapcore.ErrorLevel, &config, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})

	chainClient.ParallelTransactions(true)

	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Error deploying PLI token")

	mockETHLinkFeed, err := contractDeployer.DeployMockETHPLIFeed(ocr2vrf_constants.LinkEthFeedResponse)
	require.NoError(t, err, "Error deploying Mock ETH/PLI Feed")

	_, _, vrfBeaconContract, consumerContract, subID := ocr2vrf_actions.SetupOCR2VRFUniverse(
		t,
		linkToken,
		mockETHLinkFeed,
		contractDeployer,
		chainClient,
		nodeAddresses,
		pluginNodes,
		testNetwork,
	)

	//Request and Redeem Randomness
	requestID := ocr2vrf_actions.RequestAndRedeemRandomness(
		t,
		consumerContract,
		chainClient,
		vrfBeaconContract,
		ocr2vrf_constants.NumberOfRandomWordsToRequest,
		subID,
		ocr2vrf_constants.ConfirmationDelay,
		ocr2vrf_constants.RandomnessRedeemTransmissionEventTimeout,
	)

	for i := uint16(0); i < ocr2vrf_constants.NumberOfRandomWordsToRequest; i++ {
		randomness, err := consumerContract.GetRandomnessByRequestId(testcontext.Get(t), requestID, big.NewInt(int64(i)))
		require.NoError(t, err)
		l.Info().Interface("Random Number", randomness).Interface("Randomness Number Index", i).Msg("Randomness retrieved from Consumer contract")
		require.NotEqual(t, 0, randomness.Uint64(), "Randomness retrieved from Consumer contract give an answer other than 0")
	}
}

func TestOCR2VRFFulfillmentModel(t *testing.T) {
	t.Parallel()
	t.Skip("VRFv3 is on pause, skipping")
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig("Smoke", tc.OCR2)
	if err != nil {
		t.Fatal(err)
	}

	testEnvironment, testNetwork := setupOCR2VRFEnvironment(t)
	if testEnvironment.WillUseRemoteRunner() {
		return
	}

	chainClient, err := blockchain.NewEVMClient(testNetwork, testEnvironment, l)
	require.NoError(t, err, "Error connecting to blockchain")
	contractDeployer, err := contracts.NewContractDeployer(chainClient, l)
	require.NoError(t, err, "Error building contract deployer")
	pluginNodes, err := client.ConnectPluginNodes(testEnvironment)
	require.NoError(t, err, "Error connecting to Plugin nodes")
	nodeAddresses, err := actions.PluginNodeAddresses(pluginNodes)
	require.NoError(t, err, "Retreiving on-chain wallet addresses for plugin nodes shouldn't fail")

	t.Cleanup(func() {
		err := actions.TeardownSuite(t, testEnvironment, pluginNodes, nil, zapcore.ErrorLevel, &config, chainClient)
		require.NoError(t, err, "Error tearing down environment")
	})

	chainClient.ParallelTransactions(true)

	linkToken, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(t, err, "Error deploying PLI token")

	mockETHLinkFeed, err := contractDeployer.DeployMockETHPLIFeed(ocr2vrf_constants.LinkEthFeedResponse)
	require.NoError(t, err, "Error deploying Mock ETH/PLI Feed")

	_, _, vrfBeaconContract, consumerContract, subID := ocr2vrf_actions.SetupOCR2VRFUniverse(
		t,
		linkToken,
		mockETHLinkFeed,
		contractDeployer,
		chainClient,
		nodeAddresses,
		pluginNodes,
		testNetwork,
	)

	requestID := ocr2vrf_actions.RequestRandomnessFulfillmentAndWaitForFulfilment(
		t,
		consumerContract,
		chainClient,
		vrfBeaconContract,
		ocr2vrf_constants.NumberOfRandomWordsToRequest,
		subID,
		ocr2vrf_constants.ConfirmationDelay,
		ocr2vrf_constants.RandomnessFulfilmentTransmissionEventTimeout,
	)

	for i := uint16(0); i < ocr2vrf_constants.NumberOfRandomWordsToRequest; i++ {
		randomness, err := consumerContract.GetRandomnessByRequestId(testcontext.Get(t), requestID, big.NewInt(int64(i)))
		require.NoError(t, err, "Error getting Randomness result from Consumer Contract")
		l.Info().Interface("Random Number", randomness).Interface("Randomness Number Index", i).Msg("Randomness Fulfillment retrieved from Consumer contract")
		require.NotEqual(t, 0, randomness.Uint64(), "Randomness Fulfillment retrieved from Consumer contract give an answer other than 0")
	}
}

func setupOCR2VRFEnvironment(t *testing.T) (testEnvironment *environment.Environment, testNetwork blockchain.EVMNetwork) {
	if ocr2vrfSmokeConfig == nil {
		c, err := testconfig.GetConfig("Smoke", tc.OCR2VRF)
		if err != nil {
			t.Fatal(err)
		}
		ocr2vrfSmokeConfig = &c
	}

	testNetwork = networks.MustGetSelectedNetworkConfig(ocr2vrfSmokeConfig.Network)[0]
	evmConfig := eth.New(nil)
	if !testNetwork.Simulated {
		evmConfig = eth.New(&eth.Props{
			NetworkName: testNetwork.Name,
			Simulated:   testNetwork.Simulated,
			WsURLs:      testNetwork.URLs,
		})
	}

	var overrideFn = func(_ interface{}, target interface{}) {
		ctf_config.MustConfigOverridePluginVersion(ocr2vrfSmokeConfig.PluginImage, target)
		ctf_config.MightConfigOverridePyroscopeKey(ocr2vrfSmokeConfig.Pyroscope, target)
	}

	cd := plugin.NewWithOverride(0, map[string]interface{}{
		"replicas": 6,
		"toml": networks.AddNetworkDetailedConfig(
			config.BaseOCR2Config,
			ocr2vrfSmokeConfig.Pyroscope,
			config.DefaultOCR2VRFNetworkDetailTomlConfig,
			testNetwork,
		),
	}, ocr2vrfSmokeConfig.PluginImage, overrideFn)

	testEnvironment = environment.New(&environment.Config{
		NamespacePrefix: fmt.Sprintf("smoke-ocr2vrf-%s", strings.ReplaceAll(strings.ToLower(testNetwork.Name), " ", "-")),
		Test:            t,
	}).
		AddHelm(evmConfig).
		AddHelm(cd)
	err := testEnvironment.Run()

	require.NoError(t, err, "Error running test environment")

	return testEnvironment, testNetwork
}
