package wasm

import (
	"github.com/goplugin/plugin-common/pkg/logger"
	"github.com/goplugin/plugin-common/pkg/workflows/sdk"
)

type Runtime struct {
	fetchFn func(req sdk.FetchRequest) (sdk.FetchResponse, error)
	logger  logger.Logger
}

type RuntimeConfig struct {
	MaxFetchResponseSizeBytes int64
}

const (
	defaultMaxFetchResponseSizeBytes = 5 * 1024
)

func defaultRuntimeConfig() *RuntimeConfig {
	return &RuntimeConfig{
		MaxFetchResponseSizeBytes: defaultMaxFetchResponseSizeBytes,
	}
}

var _ sdk.Runtime = (*Runtime)(nil)

func (r *Runtime) Fetch(req sdk.FetchRequest) (sdk.FetchResponse, error) {
	return r.fetchFn(req)
}

func (r *Runtime) Logger() logger.Logger {
	return r.logger
}
