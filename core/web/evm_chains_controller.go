package web

import (
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/relay"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

var ErrEVMNotEnabled = errChainDisabled{name: "EVM", tomlKey: "EVM.Enabled"}

func NewEVMChainsController(app plugin.Application) ChainsController {
	return newChainsController[presenters.EVMChainResource](
		relay.NetworkEVM,
		app.GetRelayers().List(plugin.FilterRelayersByType(relay.NetworkEVM)),
		ErrEVMNotEnabled,
		presenters.NewEVMChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
