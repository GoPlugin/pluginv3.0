package chains

import (
	"errors"
)

var (
	// ErrChainIDEmpty is returned when chain is required but was empty.
	ErrChainIDEmpty       = errors.New("chain id empty")
	ErrNotFound           = errors.New("not found")
	ErrMultipleChainFound = errors.New("multiple chains found with the same chain ID")
)

// ChainOpts holds options for configuring a Chain
type ChainOpts interface {
	Validate() error
}
