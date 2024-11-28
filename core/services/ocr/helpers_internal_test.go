package ocr

import (
	"testing"

	"github.com/goplugin/plugin-common/pkg/sqlutil"
	"github.com/goplugin/pluginv3.0/v2/core/logger"
)

func (c *ConfigOverriderImpl) ExportedUpdateFlagsStatus() error {
	return c.updateFlagsStatus()
}

func NewTestDB(t *testing.T, ds sqlutil.DataSource, oracleSpecID int32) *db {
	return NewDB(ds, oracleSpecID, logger.TestLogger(t))
}
