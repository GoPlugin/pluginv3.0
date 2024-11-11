package web

import (
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/keystore/keys/solkey"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

func NewSolanaKeysController(app plugin.Application) KeysController {
	return NewKeysController[solkey.Key, presenters.SolanaKeyResource](app.GetKeyStore().Solana(), app.GetLogger(), app.GetAuditLogger(),
		"solanaKey", presenters.NewSolanaKeyResource, presenters.NewSolanaKeyResources)
}
