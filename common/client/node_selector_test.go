package client

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goplugin/pluginv3.0/v2/common/types"
)

func TestNodeSelector(t *testing.T) {
	// rest of the tests are located in specific node selectors tests
	t.Run("panics on unknown type", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = newNodeSelector[types.ID, RPCClient[types.ID, Head]]("unknown", nil)
		})
	})
}
