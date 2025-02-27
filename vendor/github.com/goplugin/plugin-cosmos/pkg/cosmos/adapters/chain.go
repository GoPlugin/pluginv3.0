package adapters

import (
	"github.com/goplugin/plugin-common/pkg/types"

	"github.com/goplugin/plugin-cosmos/pkg/cosmos/client"
	"github.com/goplugin/plugin-cosmos/pkg/cosmos/config"
)

type Chain interface {
	types.ChainService

	ID() string
	Config() config.Config
	TxManager() TxManager
	// Reader returns a new Reader. If nodeName is provided, the underlying client must use that node.
	Reader(nodeName string) (client.Reader, error)
}
