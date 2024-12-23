package experiments

import (
	"math/big"
	"testing"
	"time"

	seth_utils "github.com/goplugin/plugin-testing-framework/lib/utils/seth"

	"github.com/stretchr/testify/require"

	"github.com/goplugin/plugin-testing-framework/lib/logging"
	"github.com/goplugin/plugin-testing-framework/lib/networks"
	"github.com/goplugin/pluginv3.0/integration-tests/actions"
	"github.com/goplugin/pluginv3.0/integration-tests/contracts"
	tc "github.com/goplugin/pluginv3.0/integration-tests/testconfig"
)

func TestGasExperiment(t *testing.T) {
	l := logging.GetTestLogger(t)
	config, err := tc.GetConfig([]string{"Soak"}, tc.OCR)
	require.NoError(t, err, "Error getting config")

	network := networks.MustGetSelectedNetworkConfig(config.GetNetworkConfig())[0]
	seth, err := seth_utils.GetChainClient(&config, network)
	require.NoError(t, err, "Error creating seth client")

	_, err = actions.SendFunds(l, seth, actions.FundsToSendPayload{
		ToAddress:  seth.Addresses[0],
		Amount:     big.NewInt(10_000_000),
		PrivateKey: seth.PrivateKeys[0],
	})
	require.NoError(t, err, "Error sending funds")

	for i := 0; i < 1; i++ {
		_, err = contracts.DeployLinkTokenContract(l, seth)
		require.NoError(t, err, "Error deploying PLI contract")
		time.Sleep(2 * time.Second)
	}
}
