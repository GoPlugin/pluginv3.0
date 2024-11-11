package vrfv2plus

import (
	"github.com/goplugin/pluginv3.0/integration-tests/contracts"
)

type VRFV2PlusWrapperContracts struct {
	VRFV2PlusWrapper  contracts.VRFV2PlusWrapper
	LoadTestConsumers []contracts.VRFv2PlusWrapperLoadTestConsumer
}
