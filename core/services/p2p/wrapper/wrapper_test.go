package wrapper_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/goplugin/pluginv3.0/v2/core/internal/testutils"
	"github.com/goplugin/pluginv3.0/v2/core/internal/testutils/configtest"
	"github.com/goplugin/pluginv3.0/v2/core/internal/testutils/pgtest"
	"github.com/goplugin/pluginv3.0/v2/core/logger"
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/keystore/keys/p2pkey"
	ksmocks "github.com/goplugin/pluginv3.0/v2/core/services/keystore/mocks"
	"github.com/goplugin/pluginv3.0/v2/core/services/p2p/wrapper"
)

func TestPeerWrapper_CleanStartClose(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	lggr := logger.TestLogger(t)
	port := freeport.GetOne(t)
	cfg := configtest.NewGeneralConfig(t, func(c *plugin.Config, s *plugin.Secrets) {
		enabled := true
		c.Capabilities.Peering.V2.Enabled = &enabled
		c.Capabilities.Peering.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", port)}
	})
	keystoreP2P := ksmocks.NewP2P(t)
	key, err := p2pkey.NewV2()
	require.NoError(t, err)
	keystoreP2P.On("GetOrFirst", mock.Anything).Return(key, nil)

	wrapper := wrapper.NewExternalPeerWrapper(keystoreP2P, cfg.Capabilities().Peering(), db, lggr)
	require.NotNil(t, wrapper)
	require.NoError(t, wrapper.Start(testutils.Context(t)))
	require.NoError(t, wrapper.Close())
}
