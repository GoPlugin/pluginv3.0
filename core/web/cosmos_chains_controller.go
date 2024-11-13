package web

import (
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/relay"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

func NewCosmosChainsController(app plugin.Application) ChainsController {
	return newChainsController[presenters.CosmosChainResource](
		relay.NetworkCosmos,
		app.GetRelayers().List(plugin.FilterRelayersByType(relay.NetworkCosmos)),
		ErrCosmosNotEnabled,
		presenters.NewCosmosChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
