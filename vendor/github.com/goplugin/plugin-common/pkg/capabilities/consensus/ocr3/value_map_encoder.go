package ocr3

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/goplugin/plugin-common/pkg/capabilities/consensus/ocr3/types"
	"github.com/goplugin/plugin-common/pkg/values"
)

type ValueMapEncoder struct{}

func (v ValueMapEncoder) Encode(_ context.Context, input values.Map) ([]byte, error) {
	opts := proto.MarshalOptions{Deterministic: true}
	return opts.Marshal(values.Proto(&input))
}

var _ types.Encoder = (*ValueMapEncoder)(nil)
