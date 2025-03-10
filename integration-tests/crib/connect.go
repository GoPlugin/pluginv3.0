package crib

import (
	"net/http"
	"time"

	"github.com/goplugin/plugin-testing-framework/lib/crib"

	"github.com/pkg/errors"

	"github.com/goplugin/plugin-testing-framework/seth"

	"github.com/goplugin/plugin-testing-framework/lib/utils/ptr"
	seth_utils "github.com/goplugin/plugin-testing-framework/lib/utils/seth"
	tc "github.com/goplugin/pluginv3.0/integration-tests/testconfig"

	"github.com/goplugin/plugin-testing-framework/lib/blockchain"
	msClient "github.com/goplugin/plugin-testing-framework/lib/client"
	"github.com/goplugin/pluginv3.0/integration-tests/client"
)

func setSethConfig(cfg tc.TestConfig, netWSURL string, netHTTPURL string, headers http.Header) {
	netName := "CRIB_SIMULATED"
	cfg.Network.SelectedNetworks = []string{netName}
	cfg.Network.RpcHttpUrls = map[string][]string{}
	cfg.Network.RpcHttpUrls[netName] = []string{netHTTPURL}
	cfg.Network.RpcWsUrls = map[string][]string{}
	cfg.Network.RpcWsUrls[netName] = []string{netWSURL}
	cfg.Seth.EphemeralAddrs = ptr.Ptr(int64(0))
	cfg.Seth.RPCHeaders = headers
}

// ConnectRemote connects to a local environment, see https://github.com/goplugin/crib/tree/main/core
// connects to default CRIB network if simulated = true
func ConnectRemote() (
	*seth.Client,
	*msClient.MockserverClient,
	*client.PluginK8sClient,
	[]*client.PluginK8sClient,
	*crib.CoreDONConnectionConfig,
	error,
) {
	vars, err := crib.CoreDONSimulatedConnection()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	// TODO: move all the parts of ConnectRemote() to CTF when Seth config refactor is finalized
	config, err := tc.GetConfig([]string{"CRIB"}, tc.OCR)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	var sethClient *seth.Client
	switch vars.Network {
	case "geth":
		setSethConfig(config, vars.NetworkWSURL, vars.NetworkHTTPURL, vars.BlockchainNodeHeaders)
		net := blockchain.EVMNetwork{
			Name:                      vars.Network,
			Simulated:                 true,
			SupportsEIP1559:           true,
			ClientImplementation:      blockchain.EthereumClientImplementation,
			ChainID:                   vars.ChainID,
			PrivateKeys:               vars.PrivateKeys,
			URLs:                      []string{vars.NetworkWSURL},
			HTTPURLs:                  []string{vars.NetworkHTTPURL},
			PluginTransactionLimit: 500000,
			Timeout:                   blockchain.StrDuration{Duration: 2 * time.Minute},
			MinimumConfirmations:      1,
			GasEstimationBuffer:       10000,
			Headers:                   vars.BlockchainNodeHeaders,
		}
		sethClient, err = seth_utils.GetChainClient(config, net)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
	default:
		return nil, nil, nil, nil, nil, errors.New("CRIB network is not supported")
	}
	// bootstrap node
	clClients := make([]*client.PluginK8sClient, 0)
	c, err := client.NewPluginK8sClient(&client.PluginConfig{
		Email:      client.CLNodeTestEmail,
		Password:   client.CLNodeTestPassword,
		URL:        vars.NodeURLs[0],
		InternalIP: vars.NodeInternalDNS[0],
		Headers:    vars.NodeHeaders[0],
	}, vars.NodeInternalDNS[0], vars.Namespace)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	clClients = append(clClients, c)
	// all the other nodes, indices of nodes in CRIB starts with 1
	for i := 1; i < vars.Nodes; i++ {
		cl, err := client.NewPluginK8sClient(&client.PluginConfig{
			Email:      client.CLNodeTestEmail,
			Password:   client.CLNodeTestPassword,
			URL:        vars.NodeURLs[i],
			InternalIP: vars.NodeInternalDNS[i],
			Headers:    vars.NodeHeaders[i],
		}, vars.NodeInternalDNS[i], vars.Namespace)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
		clClients = append(clClients, cl)
	}
	mockServerClient := msClient.NewMockserverClient(&msClient.MockserverConfig{
		LocalURL:   vars.MockserverURL,
		ClusterURL: "http://mockserver:1080",
		Headers:    vars.MockserverHeaders,
	})

	//nolint:gosec // G602 - false positive https://github.com/securego/gosec/issues/1005
	return sethClient, mockServerClient, clClients[0], clClients[1:], vars, nil
}
