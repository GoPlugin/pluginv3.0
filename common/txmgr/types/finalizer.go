package types

import (
	"github.com/goplugin/plugin-common/pkg/services"
	"github.com/goplugin/pluginv3.0/v2/common/types"
)

type Finalizer[BLOCK_HASH types.Hashable, HEAD types.Head[BLOCK_HASH]] interface {
	// interfaces for running the underlying estimator
	services.Service
	DeliverLatestHead(head HEAD) bool
}
