package web

import (
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/relay"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

// ErrStarkNetNotEnabled is returned when Starknet.Enabled is not true.
var ErrStarkNetNotEnabled = errChainDisabled{name: "StarkNet", tomlKey: "Starknet.Enabled"}

func NewStarkNetNodesController(app plugin.Application) NodesController {
	scopedNodeStatuser := NewNetworkScopedNodeStatuser(app.GetRelayers(), relay.NetworkStarkNet)

	return newNodesController[presenters.StarkNetNodeResource](
		scopedNodeStatuser, ErrStarkNetNotEnabled, presenters.NewStarkNetNodeResource, app.GetAuditLogger())
}
