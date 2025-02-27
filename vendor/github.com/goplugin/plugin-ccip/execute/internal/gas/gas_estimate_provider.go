package gas

import "github.com/goplugin/plugin-ccip/pkg/types/ccipocr3"

type EstimateProvider interface {
	CalculateMerkleTreeGas(numRequests int) uint64
	CalculateMessageMaxGas(msg ccipocr3.Message) uint64
}
