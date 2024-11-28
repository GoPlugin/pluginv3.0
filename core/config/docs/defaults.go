package docs

import (
	"log"
	"strings"

	"github.com/goplugin/plugin-common/pkg/config"
	"github.com/goplugin/pluginv3.0/v2/core/config/toml"
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin/cfgtest"
	"github.com/goplugin/pluginv3.0/v2/core/store/dialects"
)

var (
	defaults toml.Core
)

func init() {
	if err := cfgtest.DocDefaultsOnly(strings.NewReader(coreTOML), &defaults, config.DecodeTOML); err != nil {
		log.Fatalf("Failed to initialize defaults from docs: %v", err)
	}
}

func CoreDefaults() (c toml.Core) {
	c.SetFrom(&defaults)
	c.Database.Dialect = dialects.Postgres // not user visible - overridden for tests only
	c.Tracing.Attributes = make(map[string]string)
	return
}
