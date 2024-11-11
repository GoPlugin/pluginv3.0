package docs

import (
	"testing"

	"github.com/goplugin/pluginv3.0/v2/core/services/plugin/cfgtest"
)

func TestCoreDefaults_notNil(t *testing.T) {
	cfgtest.AssertFieldsNotNil(t, CoreDefaults())
}
