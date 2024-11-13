package web

import (
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/relay"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

// ErrCosmosNotEnabled is returned when COSMOS_ENABLED is not true.
var ErrCosmosNotEnabled = errChainDisabled{name: "Cosmos", tomlKey: "Cosmos.Enabled"}

func NewCosmosNodesController(app plugin.Application) NodesController {
	scopedNodeStatuser := NewNetworkScopedNodeStatuser(app.GetRelayers(), relay.NetworkCosmos)

	return newNodesController[presenters.CosmosNodeResource](
		scopedNodeStatuser, ErrCosmosNotEnabled, presenters.NewCosmosNodeResource, app.GetAuditLogger(),
	)
}
