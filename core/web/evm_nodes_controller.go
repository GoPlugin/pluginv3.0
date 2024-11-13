package web

import (
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/relay"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

func NewEVMNodesController(app plugin.Application) NodesController {
	scopedNodeStatuser := NewNetworkScopedNodeStatuser(app.GetRelayers(), relay.NetworkEVM)

	return newNodesController[presenters.EVMNodeResource](
		scopedNodeStatuser, ErrEVMNotEnabled, presenters.NewEVMNodeResource, app.GetAuditLogger())
}
