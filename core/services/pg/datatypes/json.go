package datatypes

import (
	"github.com/goplugin/plugin-common/pkg/sqlutil"
)

// JSON defined JSON data type, need to implements driver.Valuer, sql.Scanner interface
// Deprecated: Use sqlutil.JSON instead
type JSON = sqlutil.JSON
