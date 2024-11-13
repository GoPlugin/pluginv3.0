package actions

import (
	"github.com/rs/zerolog"

	"github.com/goplugin/plugin-testing-framework/seth"

	"github.com/goplugin/pluginv3.0/integration-tests/contracts"
	tc "github.com/goplugin/pluginv3.0/integration-tests/testconfig"
)

// LinkTokenContract returns a link token contract instance. Depending on test configuration, it either deploys a new one or uses an existing one.
func LinkTokenContract(l zerolog.Logger, sethClient *seth.Client, configWithLinkToken tc.LinkTokenContractConfig) (*contracts.EthereumLinkToken, error) {
	if configWithLinkToken != nil && configWithLinkToken.UseExistingLinkTokenContract() {
		linkAddress, err := configWithLinkToken.LinkTokenContractAddress()
		if err != nil {
			return nil, err
		}

		return contracts.LoadLinkTokenContract(l, sethClient, linkAddress)
	}
	return contracts.DeployLinkTokenContract(l, sethClient)
}
