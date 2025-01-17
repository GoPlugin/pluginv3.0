package ocrbootstrap_test

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	ocrtypes "github.com/goplugin/plugin-libocr/offchainreporting2plus/types"

	"github.com/goplugin/pluginv3.0/v2/core/internal/cltest"
	"github.com/goplugin/pluginv3.0/v2/core/internal/testutils"
	"github.com/goplugin/pluginv3.0/v2/core/internal/testutils/pgtest"
	"github.com/goplugin/pluginv3.0/v2/core/logger"
	"github.com/goplugin/pluginv3.0/v2/core/services/job"
	"github.com/goplugin/pluginv3.0/v2/core/services/ocr2/testhelpers"
	"github.com/goplugin/pluginv3.0/v2/core/services/ocrbootstrap"
)

func MustInsertOCRBootstrapSpec(t *testing.T, db *sqlx.DB) job.BootstrapSpec {
	t.Helper()

	spec := job.BootstrapSpec{}
	require.NoError(t, db.Get(&spec, `INSERT INTO bootstrap_specs (
		relay, relay_config, contract_id, monitoring_endpoint, 
		blockchain_timeout, contract_config_tracker_poll_interval, contract_config_confirmations,
		created_at, updated_at) VALUES (
		'evm', '{}', $1, $2, 0, 0, 0, NOW(), NOW()
) RETURNING *`, cltest.NewEIP55Address().String(), "chain.link:1234"))
	return spec
}

func setupDB(t *testing.T) *sqlx.DB {
	t.Helper()
	return pgtest.NewSqlxDB(t)
}

func Test_DB_ReadWriteConfig(t *testing.T) {
	sqlDB := setupDB(t)

	config := ocrtypes.ContractConfig{
		ConfigDigest:          testhelpers.MakeConfigDigest(t),
		ConfigCount:           1,
		Signers:               []ocrtypes.OnchainPublicKey{{0x01}, {0x02}},
		Transmitters:          []ocrtypes.Account{"account1", "account2"},
		F:                     79,
		OnchainConfig:         []byte{0x01, 0x02},
		OffchainConfigVersion: 111,
		OffchainConfig:        []byte{0x03, 0x04},
	}
	spec := MustInsertOCRBootstrapSpec(t, sqlDB)
	lggr := logger.TestLogger(t)

	t.Run("reads and writes config", func(t *testing.T) {
		db := ocrbootstrap.NewDB(sqlDB, spec.ID, lggr)

		err := db.WriteConfig(testutils.Context(t), config)
		require.NoError(t, err)

		readConfig, err := db.ReadConfig(testutils.Context(t))
		require.NoError(t, err)

		require.Equal(t, &config, readConfig)
	})

	t.Run("updates config", func(t *testing.T) {
		db := ocrbootstrap.NewDB(sqlDB, spec.ID, lggr)

		newConfig := ocrtypes.ContractConfig{
			ConfigDigest: testhelpers.MakeConfigDigest(t),
			Signers:      []ocrtypes.OnchainPublicKey{{0x03}},
			Transmitters: []ocrtypes.Account{"test"},
		}

		err := db.WriteConfig(testutils.Context(t), newConfig)
		require.NoError(t, err)

		readConfig, err := db.ReadConfig(testutils.Context(t))
		require.NoError(t, err)

		require.Equal(t, &newConfig, readConfig)
	})

	t.Run("does not return result for wrong spec", func(t *testing.T) {
		db := ocrbootstrap.NewDB(sqlDB, spec.ID, lggr)

		err := db.WriteConfig(testutils.Context(t), config)
		require.NoError(t, err)

		db = ocrbootstrap.NewDB(sqlDB, -1, lggr)

		readConfig, err := db.ReadConfig(testutils.Context(t))
		require.NoError(t, err)

		require.Nil(t, readConfig)
	})
}
