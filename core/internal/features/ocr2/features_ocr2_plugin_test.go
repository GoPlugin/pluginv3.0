//go:build integration

package ocr2_test

import (
	"testing"

	"github.com/goplugin/pluginv3.0/v2/core/config/env"
)

func TestIntegration_OCR2_plugins(t *testing.T) {
	t.Setenv(string(env.MedianPlugin.Cmd), "plugin-feeds")
	testIntegration_OCR2(t)
}
