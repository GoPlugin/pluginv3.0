package docs_test

import (
	"strings"
	"testing"

	"github.com/kylelemons/godebug/diff"
	gotoml "github.com/pelletier/go-toml/v2"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	coscfg "github.com/goplugin/plugin-cosmos/pkg/cosmos/config"
	solcfg "github.com/goplugin/plugin-solana/pkg/solana/config"
	stkcfg "github.com/goplugin/plugin-starknet/relayer/pkg/plugin/config"

	"github.com/goplugin/plugin-common/pkg/config"

	"github.com/goplugin/pluginv3.0/v2/core/chains/evm/assets"
	"github.com/goplugin/pluginv3.0/v2/core/chains/evm/config/chaintype"
	evmcfg "github.com/goplugin/pluginv3.0/v2/core/chains/evm/config/toml"
	"github.com/goplugin/pluginv3.0/v2/core/chains/evm/types"
	"github.com/goplugin/pluginv3.0/v2/core/config/docs"
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin/cfgtest"
)

func TestDoc(t *testing.T) {
	d := gotoml.NewDecoder(strings.NewReader(docs.DocsTOML))
	d.DisallowUnknownFields() // Ensure no extra fields
	var c plugin.Config
	err := d.Decode(&c)
	var strict *gotoml.StrictMissingError
	if err != nil && strings.Contains(err.Error(), "undecoded keys: ") {
		t.Errorf("Docs contain extra fields: %v", err)
	} else if pkgerrors.As(err, &strict) {
		t.Fatal("StrictMissingError:", strict.String())
	} else {
		require.NoError(t, err)
	}

	cfgtest.AssertFieldsNotNil(t, c)

	var defaults plugin.Config
	require.NoError(t, cfgtest.DocDefaultsOnly(strings.NewReader(docs.DocsTOML), &defaults, config.DecodeTOML))

	t.Run("EVM", func(t *testing.T) {
		fallbackDefaults := evmcfg.Defaults(nil)
		docDefaults := defaults.EVM[0].Chain

		require.Equal(t, chaintype.ChainType(""), docDefaults.ChainType.ChainType())
		docDefaults.ChainType = nil

		// clean up KeySpecific as a special case
		require.Equal(t, 1, len(docDefaults.KeySpecific))
		ks := evmcfg.KeySpecific{Key: new(types.EIP55Address),
			GasEstimator: evmcfg.KeySpecificGasEstimator{PriceMax: new(assets.Wei)}}
		require.Equal(t, ks, docDefaults.KeySpecific[0])
		docDefaults.KeySpecific = nil

		// EVM.GasEstimator.BumpTxDepth doesn't have a constant default - it is derived from another field
		require.Zero(t, *docDefaults.GasEstimator.BumpTxDepth)
		docDefaults.GasEstimator.BumpTxDepth = nil

		// per-job limits are nilable
		require.Zero(t, *docDefaults.GasEstimator.LimitJobType.OCR)
		require.Zero(t, *docDefaults.GasEstimator.LimitJobType.OCR2)
		require.Zero(t, *docDefaults.GasEstimator.LimitJobType.DR)
		require.Zero(t, *docDefaults.GasEstimator.LimitJobType.Keeper)
		require.Zero(t, *docDefaults.GasEstimator.LimitJobType.VRF)
		require.Zero(t, *docDefaults.GasEstimator.LimitJobType.FM)
		docDefaults.GasEstimator.LimitJobType = evmcfg.GasLimitJobType{}

		// EIP1559FeeCapBufferBlocks doesn't have a constant default - it is derived from another field
		require.Zero(t, *docDefaults.GasEstimator.BlockHistory.EIP1559FeeCapBufferBlocks)
		docDefaults.GasEstimator.BlockHistory.EIP1559FeeCapBufferBlocks = nil

		// addresses w/o global values
		require.Zero(t, *docDefaults.FlagsContractAddress)
		require.Zero(t, *docDefaults.LinkContractAddress)
		require.Zero(t, *docDefaults.OperatorFactoryAddress)
		docDefaults.FlagsContractAddress = nil
		docDefaults.LinkContractAddress = nil
		docDefaults.OperatorFactoryAddress = nil
		require.Empty(t, docDefaults.Workflow.FromAddress)
		require.Empty(t, docDefaults.Workflow.ForwarderAddress)
		gasLimitDefault := uint64(400_000)
		require.Equal(t, &gasLimitDefault, docDefaults.Workflow.GasLimitDefault)

		docDefaults.Workflow.FromAddress = nil
		docDefaults.Workflow.ForwarderAddress = nil
		docDefaults.Workflow.GasLimitDefault = &gasLimitDefault
		docDefaults.NodePool.Errors = evmcfg.ClientErrors{}

		// Transactions.AutoPurge configs are only set if the feature is enabled
		docDefaults.Transactions.AutoPurge.DetectionApiUrl = nil
		docDefaults.Transactions.AutoPurge.Threshold = nil
		docDefaults.Transactions.AutoPurge.MinAttempts = nil

		// Fallback DA oracle is not set
		docDefaults.GasEstimator.DAOracle = evmcfg.DAOracle{}

		assertTOML(t, fallbackDefaults, docDefaults)
	})

	t.Run("Cosmos", func(t *testing.T) {
		var fallbackDefaults coscfg.TOMLConfig
		fallbackDefaults.SetDefaults()

		assertTOML(t, fallbackDefaults.Chain, defaults.Cosmos[0].Chain)
	})

	t.Run("Solana", func(t *testing.T) {
		var fallbackDefaults solcfg.TOMLConfig
		fallbackDefaults.SetDefaults()

		assertTOML(t, fallbackDefaults.Chain, defaults.Solana[0].Chain)
	})

	t.Run("Starknet", func(t *testing.T) {
		var fallbackDefaults stkcfg.TOMLConfig
		fallbackDefaults.SetDefaults()

		assertTOML(t, fallbackDefaults.Chain, defaults.Starknet[0].Chain)
	})
}

func assertTOML[T any](t *testing.T, fallback, docs T) {
	t.Helper()
	t.Logf("fallback: %#v", fallback)
	t.Logf("docs: %#v", docs)
	fb, err := gotoml.Marshal(fallback)
	require.NoError(t, err)
	db, err := gotoml.Marshal(docs)
	require.NoError(t, err)
	fs, ds := string(fb), string(db)
	assert.Equal(t, fs, ds, diff.Diff(fs, ds))
}
