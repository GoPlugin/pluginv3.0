package txm

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"golang.org/x/crypto/ripemd160" //nolint: staticcheck

	"github.com/goplugin/plugin-common/pkg/loop"
)

type accountInfo struct {
	Account string
	PubKey  *secp256k1.PubKey
}

// keystoreAdapter adapts a Cosmos loop.Keystore to translate public keys into bech32-prefixed account addresses.
type keystoreAdapter struct {
	keystore        loop.Keystore
	accountPrefix   string
	mutex           sync.RWMutex
	addressToPubKey map[string]*accountInfo
}

func newKeystoreAdapter(keystore loop.Keystore, accountPrefix string) *keystoreAdapter {
	return &keystoreAdapter{
		keystore:        keystore,
		accountPrefix:   accountPrefix,
		addressToPubKey: make(map[string]*accountInfo),
	}
}

func (ka *keystoreAdapter) updateMappingLocked(ctx context.Context) error {
	accounts, err := ka.keystore.Accounts(ctx)
	if err != nil {
		return err
	}

	// similar to cosmos-sdk, cache and re-use calculated bech32 addresses to prevent duplicated work.
	// ref: https://github.com/cosmos/cosmos-sdk/blob/3b509c187e1643757f5ef8a0b5ae3decca0c7719/types/address.go#L705

	type cacheEntry struct {
		bech32Addr  string
		accountInfo *accountInfo
	}
	accountCache := make(map[string]cacheEntry, len(ka.addressToPubKey))
	for bech32Addr, accountInfo := range ka.addressToPubKey {
		accountCache[accountInfo.Account] = cacheEntry{bech32Addr: bech32Addr, accountInfo: accountInfo}
	}

	addressToPubKey := make(map[string]*accountInfo, len(accounts))
	for _, account := range accounts {
		if prevEntry, ok := accountCache[account]; ok {
			addressToPubKey[prevEntry.bech32Addr] = prevEntry.accountInfo
			continue
		}
		pubKeyBytes, err := hex.DecodeString(account)
		if err != nil {
			return err
		}

		if len(pubKeyBytes) != secp256k1.PubKeySize {
			return errors.New("length of pubkey is incorrect")
		}

		sha := sha256.Sum256(pubKeyBytes)
		hasherRIPEMD160 := ripemd160.New()
		_, _ = hasherRIPEMD160.Write(sha[:])
		address := crypto.Address(hasherRIPEMD160.Sum(nil))

		bech32Addr, err := bech32.ConvertAndEncode(ka.accountPrefix, address)
		if err != nil {
			return err
		}

		addressToPubKey[bech32Addr] = &accountInfo{
			Account: account,
			PubKey:  &secp256k1.PubKey{Key: pubKeyBytes},
		}
	}

	ka.addressToPubKey = addressToPubKey
	return nil
}

func (ka *keystoreAdapter) lookup(ctx context.Context, id string) (*accountInfo, error) {
	ka.mutex.RLock()
	ai, ok := ka.addressToPubKey[id]
	ka.mutex.RUnlock()
	if !ok {
		// try updating the mapping once, incase there was an update on the keystore.
		ka.mutex.Lock()
		err := ka.updateMappingLocked(ctx)
		if err != nil {
			ka.mutex.Unlock()
			return nil, err
		}
		ai, ok = ka.addressToPubKey[id]
		ka.mutex.Unlock()
		if !ok {
			return nil, errors.New("No such id")
		}
	}
	return ai, nil
}

func (ka *keystoreAdapter) Sign(ctx context.Context, id string, hash []byte) ([]byte, error) {
	accountInfo, err := ka.lookup(ctx, id)
	if err != nil {
		return nil, err
	}
	return ka.keystore.Sign(ctx, accountInfo.Account, hash)
}

// Returns the cosmos PubKey associated with the prefixed address.
func (ka *keystoreAdapter) PubKey(ctx context.Context, address string) (cryptotypes.PubKey, error) {
	accountInfo, err := ka.lookup(ctx, address)
	if err != nil {
		return nil, err
	}
	return accountInfo.PubKey, nil
}
