package web

import (
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/relay"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

func NewStarkNetChainsController(app plugin.Application) ChainsController {
	return newChainsController(
		relay.NetworkStarkNet,
		app.GetRelayers().List(plugin.FilterRelayersByType(relay.NetworkStarkNet)),
		ErrStarkNetNotEnabled,
		presenters.NewStarkNetChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
