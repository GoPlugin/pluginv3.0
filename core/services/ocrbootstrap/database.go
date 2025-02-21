package ocrbootstrap

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	ocrtypes "github.com/goplugin/plugin-libocr/offchainreporting2plus/types"

	"github.com/goplugin/plugin-common/pkg/sqlutil"
	"github.com/goplugin/pluginv3.0/v2/core/logger"
)

type db struct {
	ds           sqlutil.DataSource
	oracleSpecID int32
	lggr         logger.Logger
}

var _ ocrtypes.ConfigDatabase = &db{}

// NewDB returns a new DB scoped to this oracleSpecID
func NewDB(ds sqlutil.DataSource, bootstrapSpecID int32, lggr logger.Logger) *db {
	return &db{ds, bootstrapSpecID, lggr}
}

func (d *db) ReadConfig(ctx context.Context) (c *ocrtypes.ContractConfig, err error) {
	q := d.ds.QueryRowxContext(ctx, `
SELECT
	config_digest,
	config_count,
	signers,
	transmitters,
	f,
	onchain_config,
	offchain_config_version,
	offchain_config
FROM bootstrap_contract_configs
WHERE bootstrap_spec_id = $1
LIMIT 1`, d.oracleSpecID)

	c = new(ocrtypes.ContractConfig)

	digest := []byte{}
	signers := [][]byte{}
	transmitters := [][]byte{}

	err = q.Scan(
		&digest,
		&c.ConfigCount,
		(*pq.ByteaArray)(&signers),
		(*pq.ByteaArray)(&transmitters),
		&c.F,
		&c.OnchainConfig,
		&c.OffchainConfigVersion,
		&c.OffchainConfig,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "ReadConfig failed")
	}

	copy(c.ConfigDigest[:], digest)

	c.Signers = []ocrtypes.OnchainPublicKey{}
	for _, s := range signers {
		c.Signers = append(c.Signers, s)
	}

	c.Transmitters = []ocrtypes.Account{}
	for _, t := range transmitters {
		transmitter := ocrtypes.Account(t)
		c.Transmitters = append(c.Transmitters, transmitter)
	}

	return
}

func (d *db) WriteConfig(ctx context.Context, c ocrtypes.ContractConfig) error {
	var signers [][]byte
	for _, s := range c.Signers {
		signers = append(signers, []byte(s))
	}
	_, err := d.ds.ExecContext(ctx, `
INSERT INTO bootstrap_contract_configs (
	bootstrap_spec_id,
	config_digest,
	config_count,
	signers,
	transmitters,
	f,
	onchain_config,
	offchain_config_version,
	offchain_config,
	created_at,
	updated_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
ON CONFLICT (bootstrap_spec_id) DO UPDATE SET
	config_digest = EXCLUDED.config_digest,
	config_count = EXCLUDED.config_count,
	signers = EXCLUDED.signers,
	transmitters = EXCLUDED.transmitters,
	f = EXCLUDED.f,
	onchain_config = EXCLUDED.onchain_config,
	offchain_config_version = EXCLUDED.offchain_config_version,
	offchain_config = EXCLUDED.offchain_config,
	updated_at = NOW()
`,
		d.oracleSpecID,
		c.ConfigDigest,
		c.ConfigCount,
		pq.ByteaArray(signers),
		c.Transmitters,
		c.F,
		c.OnchainConfig,
		c.OffchainConfigVersion,
		c.OffchainConfig,
	)

	return errors.Wrap(err, "WriteConfig failed")
}
