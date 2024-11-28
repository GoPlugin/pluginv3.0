package evm

import (
	"github.com/goplugin/plugin-common/pkg/loop"

	"github.com/goplugin/pluginv3.0/v2/core/chains/legacyevm"
	"github.com/goplugin/pluginv3.0/v2/core/services/relay"
)

//go:generate mockery --quiet --name LoopRelayAdapter --output ./mocks/ --case=underscore
type LoopRelayAdapter interface {
	loop.Relayer
	Chain() legacyevm.Chain
}
type LoopRelayer struct {
	loop.Relayer
	ext EVMChainRelayerExtender
}

var _ loop.Relayer = &LoopRelayer{}

func NewLoopRelayServerAdapter(r *Relayer, cs EVMChainRelayerExtender) *LoopRelayer {
	ra := relay.NewServerAdapter(r, cs)
	return &LoopRelayer{
		Relayer: ra,
		ext:     cs,
	}
}

func (la *LoopRelayer) Chain() legacyevm.Chain {
	return la.ext.Chain()
}
