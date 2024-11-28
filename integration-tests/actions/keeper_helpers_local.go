package actions

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/goplugin/pluginv3.0/integration-tests/client"
	"github.com/goplugin/pluginv3.0/integration-tests/contracts"
)

func CreateKeeperJobsLocal(
	l zerolog.Logger,
	pluginNodes []*client.PluginClient,
	keeperRegistry contracts.KeeperRegistry,
	ocrConfig contracts.OCRv2Config,
	evmChainID string,
) ([]*client.Job, error) {
	// Send keeper jobs to registry and plugin nodes
	primaryNode := pluginNodes[0]
	primaryNodeAddress, err := primaryNode.PrimaryEthAddress()
	if err != nil {
		l.Error().Err(err).Msg("Reading ETH Keys from Plugin Client shouldn't fail")
		return nil, err
	}
	nodeAddresses, err := PluginNodeAddressesLocal(pluginNodes)
	if err != nil {
		l.Error().Err(err).Msg("Retrieving on-chain wallet addresses for plugin nodes shouldn't fail")
		return nil, err
	}
	nodeAddressesStr, payees := make([]string, 0), make([]string, 0)
	for _, cla := range nodeAddresses {
		nodeAddressesStr = append(nodeAddressesStr, cla.Hex())
		payees = append(payees, primaryNodeAddress)
	}
	err = keeperRegistry.SetKeepers(nodeAddressesStr, payees, ocrConfig)
	if err != nil {
		l.Error().Err(err).Msg("Setting keepers in the registry shouldn't fail")
		return nil, err
	}
	jobs := []*client.Job{}
	for _, pluginNode := range pluginNodes {
		pluginNodeAddress, err := pluginNode.PrimaryEthAddress()
		if err != nil {
			l.Error().Err(err).Msg("Error retrieving plugin node address")
			return nil, err
		}
		job, err := pluginNode.MustCreateJob(&client.KeeperJobSpec{
			Name:                     fmt.Sprintf("keeper-test-%s", keeperRegistry.Address()),
			ContractAddress:          keeperRegistry.Address(),
			FromAddress:              pluginNodeAddress,
			EVMChainID:               evmChainID,
			MinIncomingConfirmations: 1,
		})
		if err != nil {
			l.Error().Err(err).Msg("Creating KeeperV2 Job shouldn't fail")
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}
