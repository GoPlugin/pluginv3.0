package changeset

import (
	"github.com/goplugin/pluginv3.0/integration-tests/deployment"
	llodeployment "github.com/goplugin/pluginv3.0/integration-tests/deployment/llo"
)

func DeployChannelConfigStoreChangeSet(env deployment.Environment, c llodeployment.DeployLLOContractConfig) (deployment.ChangesetOutput, error) {
	ab := deployment.NewMemoryAddressBook()
	err := llodeployment.DeployChannelConfigStore(env, ab, c)
	if err != nil {
		env.Logger.Errorw("Failed to deploy ChannelConfigStore", "err", err, "addresses", ab)
		return deployment.ChangesetOutput{AddressBook: ab}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}
