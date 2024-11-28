package web

import (
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/relay"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

func NewSolanaChainsController(app plugin.Application) ChainsController {
	return newChainsController(
		relay.NetworkSolana,
		app.GetRelayers().List(plugin.FilterRelayersByType(relay.NetworkSolana)),
		ErrSolanaNotEnabled,
		presenters.NewSolanaChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
