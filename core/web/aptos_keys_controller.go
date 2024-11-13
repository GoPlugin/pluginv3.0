package web

import (
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/keystore/keys/aptoskey"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

func NewAptosKeysController(app plugin.Application) KeysController {
	return NewKeysController[aptoskey.Key, presenters.AptosKeyResource](app.GetKeyStore().Aptos(), app.GetLogger(), app.GetAuditLogger(),
		"aptosKey", presenters.NewAptosKeyResource, presenters.NewAptosKeyResources)
}
