package job

import (
	"context"

	"github.com/goplugin/plugin-common/pkg/workflows/sdk"
)

type WorkflowSpecFactory interface {
	Spec(ctx context.Context, workflow, config string) (sdk.WorkflowSpec, []byte, string, error)
	RawSpec(ctx context.Context, workflow, config string) ([]byte, error)
}

var workflowSpecFactories = map[WorkflowSpecType]WorkflowSpecFactory{
	YamlSpec:        YAMLSpecFactory{},
	WASMFile:        WasmFileSpecFactory{},
	DefaultSpecType: YAMLSpecFactory{},
}
