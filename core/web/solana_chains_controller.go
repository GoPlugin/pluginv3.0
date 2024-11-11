package web

import (
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/relay"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

func NewSolanaChainsController(app plugin.Application) ChainsController {
	return newChainsController(
		relay.Solana,
		app.GetRelayers().List(plugin.FilterRelayersByType(relay.Solana)),
		ErrSolanaNotEnabled,
		presenters.NewSolanaChainResource,
		app.GetLogger(),
		app.GetAuditLogger())
}
