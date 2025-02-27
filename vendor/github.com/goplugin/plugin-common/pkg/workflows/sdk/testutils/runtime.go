package testutils

import (
	"github.com/goplugin/plugin-common/pkg/logger"
	"github.com/goplugin/plugin-common/pkg/workflows/sdk"
)

type NoopRuntime struct{}

var _ sdk.Runtime = &NoopRuntime{}

func (nr *NoopRuntime) Fetch(sdk.FetchRequest) (sdk.FetchResponse, error) {
	return sdk.FetchResponse{}, nil
}

func (nr *NoopRuntime) Logger() logger.Logger {
	l, _ := logger.New()
	return l
}
