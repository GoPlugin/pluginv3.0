package vrfv1

import (
	"fmt"

	"github.com/goplugin/plugin-testing-framework/blockchain"
	"github.com/goplugin/pluginv3.0/integration-tests/contracts"
)

var (
	ErrDeployBHSV1            = "error deploying BlockHashStoreV1 contract"
	ErrDeployVRFCootrinatorV1 = "error deploying VRFv1 Coordinator contract"
	ErrDeployVRFConsumerV1    = "error deploying VRFv1 Consumer contract"
)

type Contracts struct {
	BHS         contracts.BlockHashStore
	Coordinator contracts.VRFCoordinator
	Consumer    contracts.VRFConsumer
}

func DeployVRFContracts(cd contracts.ContractDeployer, bc blockchain.EVMClient, lt contracts.LinkToken) (*Contracts, error) {
	bhs, err := cd.DeployBlockhashStore()
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployBHSV1, err)
	}
	coordinator, err := cd.DeployVRFCoordinator(lt.Address(), bhs.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployVRFCootrinatorV1, err)
	}
	consumer, err := cd.DeployVRFConsumer(lt.Address(), coordinator.Address())
	if err != nil {
		return nil, fmt.Errorf("%s, err %w", ErrDeployVRFConsumerV1, err)
	}
	if err := bc.WaitForEvents(); err != nil {
		return nil, err
	}
	return &Contracts{bhs, coordinator, consumer}, nil
}
