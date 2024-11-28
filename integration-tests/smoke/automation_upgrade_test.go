package smoke

import (
	"testing"

	tc "github.com/goplugin/pluginv3.0/integration-tests/testconfig"
)

func TestAutomationNodeUpgrade(t *testing.T) {
	config, err := tc.GetConfig(t.Name(), tc.Automation)
	if err != nil {
		t.Fatal(err, "Error getting config")
	}
	SetupAutomationBasic(t, true, &config)
}
