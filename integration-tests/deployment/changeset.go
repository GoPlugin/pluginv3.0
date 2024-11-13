package deployment

import (
	"github.com/goplugin/ccip-owner-contracts/pkg/proposal/timelock"
)

// Services as input to CI/Async tasks
type ChangesetOutput struct {
	JobSpecs    map[string][]string
	Proposals   []timelock.MCMSWithTimelockProposal
	AddressBook AddressBook
}
