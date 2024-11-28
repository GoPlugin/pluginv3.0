package testsetups

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/goplugin/plugin-testing-framework/blockchain"
	ctfClient "github.com/goplugin/plugin-testing-framework/client"
	e "github.com/goplugin/plugin-testing-framework/k8s/environment"
	"github.com/goplugin/plugin-testing-framework/k8s/pkg/helm/plugin"
	"github.com/goplugin/plugin-testing-framework/k8s/pkg/helm/ethereum"
	"github.com/goplugin/plugin-testing-framework/k8s/pkg/helm/mockserver"
	mockservercfg "github.com/goplugin/plugin-testing-framework/k8s/pkg/helm/mockserver-cfg"

	"github.com/goplugin/pluginv3.0/integration-tests/client"
	"github.com/goplugin/pluginv3.0/integration-tests/contracts"
)

type DonChain struct {
	conf              *DonChainConfig
	EVMClient         blockchain.EVMClient
	EVMNetwork        *blockchain.EVMNetwork
	ContractDeployer  contracts.ContractDeployer
	LinkTokenContract contracts.LinkToken
	PluginNodes    []*client.PluginK8sClient
	Mockserver        *ctfClient.MockserverClient
	l                 zerolog.Logger
}

type DonChainConfig struct {
	T               *testing.T
	Env             *e.Environment
	EVMNetwork      *blockchain.EVMNetwork
	EthereumProps   *ethereum.Props
	PluginValues map[string]interface{}
}

func NewDonChain(conf *DonChainConfig, logger zerolog.Logger) *DonChain {
	return &DonChain{
		conf:       conf,
		EVMNetwork: conf.EVMNetwork,
		l:          logger,
	}
}

func (s *DonChain) Deploy() {
	var err error

	s.conf.Env.AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(s.conf.EthereumProps)).
		AddHelm(plugin.New(0, s.conf.PluginValues))

	err = s.conf.Env.Run()
	require.NoError(s.conf.T, err)

	s.initializeClients()
}

func (s *DonChain) initializeClients() {
	var err error
	network := *s.conf.EVMNetwork
	s.EVMClient, err = blockchain.NewEVMClient(network, s.conf.Env, s.l)
	require.NoError(s.conf.T, err, "Connecting to blockchain nodes shouldn't fail")

	s.ContractDeployer, err = contracts.NewContractDeployer(s.EVMClient, s.l)
	require.NoError(s.conf.T, err)

	s.PluginNodes, err = client.ConnectPluginNodes(s.conf.Env)
	require.NoError(s.conf.T, err, "Connecting to plugin nodes shouldn't fail")

	s.Mockserver, err = ctfClient.ConnectMockServer(s.conf.Env)
	require.NoError(s.conf.T, err, "Creating mockserver clients shouldn't fail")

	s.EVMClient.ParallelTransactions(true)

	s.LinkTokenContract, err = s.ContractDeployer.DeployLinkTokenContract()
	require.NoError(s.conf.T, err, "Deploying Link Token Contract shouldn't fail")
}
