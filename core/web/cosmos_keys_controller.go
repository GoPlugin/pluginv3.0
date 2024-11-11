package web

import (
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/keystore/keys/cosmoskey"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

func NewCosmosKeysController(app plugin.Application) KeysController {
	return NewKeysController[cosmoskey.Key, presenters.CosmosKeyResource](app.GetKeyStore().Cosmos(), app.GetLogger(), app.GetAuditLogger(),
		"cosmosKey", presenters.NewCosmosKeyResource, presenters.NewCosmosKeyResources)
}
