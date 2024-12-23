package types

import (
	"context"

	"github.com/goplugin/pluginv3.0/v2/common/types"
)

// KeyStore encompasses the subset of keystore used by txmgr
type KeyStore[
	// Account Address type.
	ADDR types.Hashable,
	// Chain ID type
	CHAIN_ID types.ID,
	// Chain's sequence type. For example, EVM chains use nonce, bitcoin uses UTXO.
	SEQ types.Sequence,
] interface {
	CheckEnabled(ctx context.Context, address ADDR, chainID CHAIN_ID) error
	EnabledAddressesForChain(ctx context.Context, chainId CHAIN_ID) ([]ADDR, error)
	SubscribeToKeyChanges(ctx context.Context) (ch chan struct{}, unsub func())
}
