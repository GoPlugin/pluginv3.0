//go:build integration

package ocr2_test

import (
	"testing"

	"github.com/goplugin/pluginv3.0/v2/core/config/env"
	"github.com/goplugin/pluginv3.0/v2/core/internal/testutils"
)

func TestIntegration_OCR2_plugins(t *testing.T) {
	t.Setenv(string(env.MedianPlugin.Cmd), "plugin-feeds")
	testutils.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/BCF-3417")
	testIntegration_OCR2(t)
}
