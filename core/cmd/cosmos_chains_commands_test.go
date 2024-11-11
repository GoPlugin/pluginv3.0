package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	coscfg "github.com/goplugin/plugin-cosmos/pkg/cosmos/config"

	"github.com/goplugin/pluginv3.0/v2/core/cmd"
	"github.com/goplugin/pluginv3.0/v2/core/internal/cltest"
	"github.com/goplugin/pluginv3.0/v2/core/internal/testutils/cosmostest"
)

func TestShell_IndexCosmosChains(t *testing.T) {
	t.Parallel()

	chainID := cosmostest.RandomChainID()
	chain := coscfg.TOMLConfig{
		ChainID: ptr(chainID),
		Enabled: ptr(true),
	}
	app := cosmosStartNewApplication(t, &chain)
	client, r := app.NewShellAndRenderer()

	require.Nil(t, cmd.CosmosChainClient(client).IndexChains(cltest.EmptyCLIContext()))
	chains := *r.Renders[0].(*cmd.CosmosChainPresenters)
	require.Len(t, chains, 1)
	c := chains[0]
	assert.Equal(t, chainID, c.ID)
	assertTableRenders(t, r)
}
