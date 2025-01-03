package fluxmonitorv2

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/goplugin/pluginv3.0/v2/core/services/keystore"
	"github.com/goplugin/pluginv3.0/v2/core/services/keystore/keys/ethkey"
)

// KeyStoreInterface defines an interface to interact with the keystore
type KeyStoreInterface interface {
	EnabledKeysForChain(ctx context.Context, chainID *big.Int) ([]ethkey.KeyV2, error)
	GetRoundRobinAddress(ctx context.Context, chainID *big.Int, addrs ...common.Address) (common.Address, error)
}

// KeyStore implements KeyStoreInterface
type KeyStore struct {
	keystore.Eth
}

// NewKeyStore initializes a new keystore
func NewKeyStore(ks keystore.Eth) *KeyStore {
	return &KeyStore{ks}
}
