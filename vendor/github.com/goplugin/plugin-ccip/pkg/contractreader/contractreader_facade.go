package contractreader

import (
	"context"

	"github.com/goplugin/plugin-common/pkg/types"
	"github.com/goplugin/plugin-common/pkg/types/query"
	"github.com/goplugin/plugin-common/pkg/types/query/primitives"
)

// ContractReaderFacade wraps the public functions of ContractReader in plugin-common so that we can mock it.
//
//nolint:lll // don't read this interface.
type ContractReaderFacade interface {
	GetLatestValue(ctx context.Context, readIdentifier string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error
	BatchGetLatestValues(ctx context.Context, request types.BatchGetLatestValuesRequest) (types.BatchGetLatestValuesResult, error)
	Bind(ctx context.Context, bindings []types.BoundContract) error
	Unbind(ctx context.Context, bindings []types.BoundContract) error
	QueryKey(ctx context.Context, contract types.BoundContract, filter query.KeyFilter, limitAndSort query.LimitAndSort, sequenceDataType any) ([]types.Sequence, error)
	//mustEmbedUnimplementedContractReaderServer()
}
