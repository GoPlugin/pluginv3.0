package ocr

import (
	"math/big"
	"math/rand"
	"time"

	"github.com/rs/zerolog"

	"github.com/goplugin/plugin-testing-framework/blockchain"

	client2 "github.com/goplugin/plugin-testing-framework/client"
	"github.com/goplugin/pluginv3.0/integration-tests/actions"
	"github.com/goplugin/pluginv3.0/integration-tests/client"
	"github.com/goplugin/pluginv3.0/integration-tests/contracts"
)

func SetupCluster(
	cc blockchain.EVMClient,
	cd contracts.ContractDeployer,
	workerNodes []*client.PluginK8sClient,
) (contracts.LinkToken, error) {
	err := actions.FundPluginNodes(workerNodes, cc, big.NewFloat(3))
	if err != nil {
		return nil, err
	}
	lt, err := cd.DeployLinkTokenContract()
	if err != nil {
		return nil, err
	}
	return lt, nil
}

func SetupFeed(
	cc blockchain.EVMClient,
	msClient *client2.MockserverClient,
	cd contracts.ContractDeployer,
	bootstrapNode *client.PluginK8sClient,
	workerNodes []*client.PluginK8sClient,
	lt contracts.LinkToken,
) ([]contracts.OffchainAggregator, error) {
	ocrInstances, err := actions.DeployOCRContracts(1, lt, cd, workerNodes, cc)
	if err != nil {
		return nil, err
	}
	err = actions.CreateOCRJobs(ocrInstances, bootstrapNode, workerNodes, 5, msClient, cc.GetChainID().String())
	if err != nil {
		return nil, err
	}
	return ocrInstances, nil
}

func SimulateEAActivity(
	l zerolog.Logger,
	eaChangeInterval time.Duration,
	ocrInstances []contracts.OffchainAggregator,
	workerNodes []*client.PluginK8sClient,
	msClient *client2.MockserverClient,
) {
	go func() {
		for {
			time.Sleep(eaChangeInterval)
			if err := actions.SetAllAdapterResponsesToTheSameValue(rand.Intn(1000), ocrInstances, workerNodes, msClient); err != nil {
				l.Error().Err(err).Msg("failed to update mockserver responses")
			}
		}
	}()
}
