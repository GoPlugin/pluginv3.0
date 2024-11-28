package web

import (
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/keystore/keys/dkgencryptkey"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

func NewDKGEncryptKeysController(app plugin.Application) KeysController {
	return NewKeysController[dkgencryptkey.Key, presenters.DKGEncryptKeyResource](
		app.GetKeyStore().DKGEncrypt(),
		app.GetLogger(),
		app.GetAuditLogger(),
		"dkgencryptKey",
		presenters.NewDKGEncryptKeyResource,
		presenters.NewDKGEncryptKeyResources)
}
