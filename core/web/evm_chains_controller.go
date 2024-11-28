package web

import (
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/relay"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

var ErrEVMNotEnabled = errChainDisabled{name: "EVM", tomlKey: "EVM.Enabled"}

func NewEVMChainsController(app plugin.Application) ChainsController {
	return newChainsController[presenters.EVMChainResource](
		relay.EVM,
		app.GetRelayers().List(plugin.FilterRelayersByType(relay.EVM)),
		ErrEVMNotEnabled,
		presenters.NewEVMChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
