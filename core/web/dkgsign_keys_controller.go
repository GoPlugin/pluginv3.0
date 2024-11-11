package web

import (
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/keystore/keys/dkgsignkey"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

func NewDKGSignKeysController(app plugin.Application) KeysController {
	return NewKeysController[dkgsignkey.Key, presenters.DKGSignKeyResource](
		app.GetKeyStore().DKGSign(),
		app.GetLogger(),
		app.GetAuditLogger(),
		"dkgsignKey",
		presenters.NewDKGSignKeyResource,
		presenters.NewDKGSignKeyResources)
}
