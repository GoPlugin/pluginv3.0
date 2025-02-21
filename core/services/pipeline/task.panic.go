package pipeline

import (
	"context"

	"github.com/goplugin/pluginv3.0/v2/core/logger"
)

type PanicTask struct {
	BaseTask `mapstructure:",squash"`
	Msg      string
}

var _ Task = (*PanicTask)(nil)

func (t *PanicTask) Type() TaskType {
	return TaskTypePanic
}

func (t *PanicTask) Run(_ context.Context, _ logger.Logger, vars Vars, _ []Result) (result Result, runInfo RunInfo) {
	panic(t.Msg)
}
