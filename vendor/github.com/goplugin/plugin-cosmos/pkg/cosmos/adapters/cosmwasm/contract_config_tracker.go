package cosmwasm

import (
	"context"

	"github.com/goplugin/plugin-libocr/offchainreporting2/types"

	"github.com/goplugin/plugin-cosmos/pkg/cosmos/client"
)

var _ types.ContractConfigTracker = (*ContractTracker)(nil)

type ContractTracker struct {
	*ContractCache
	chainReader client.Reader
}

func NewContractTracker(chainReader client.Reader, contract *ContractCache) *ContractTracker {
	return &ContractTracker{
		ContractCache: contract,
		chainReader:   chainReader,
	}
}

// Unused, libocr will use polling
func (ct *ContractTracker) Notify() <-chan struct{} {
	return nil
}

// TODO: seems heavy to fetch whole block rather than rpc.Status() -> SyncInfo.LatestBlockHeight
// LatestBlockHeight returns the height of the most recent block in the chain.
func (ct *ContractTracker) LatestBlockHeight(ctx context.Context) (blockHeight uint64, err error) {
	b, err := ct.chainReader.LatestBlock(ctx)
	if err != nil {
		return 0, err
	}
	return uint64(b.SdkBlock.Header.Height), nil
}
