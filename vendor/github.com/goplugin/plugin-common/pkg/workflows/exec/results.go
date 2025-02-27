package exec

import "github.com/goplugin/plugin-common/pkg/values"

type Results interface {
	ResultForStep(string) (*Result, bool)
}

type Result struct {
	Inputs  values.Value
	Outputs values.Value
	Error   error
}
