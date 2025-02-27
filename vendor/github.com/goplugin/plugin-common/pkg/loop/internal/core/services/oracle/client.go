package oracle

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	oraclepb "github.com/goplugin/plugin-common/pkg/loop/internal/pb/oracle"
	"github.com/goplugin/plugin-common/pkg/types/core"
)

var _ core.Oracle = (*client)(nil)

type client struct {
	grpc oraclepb.OracleClient
}

func NewClient(cc grpc.ClientConnInterface) *client {
	return &client{grpc: oraclepb.NewOracleClient(cc)}
}

func (c *client) Close(ctx context.Context) error {
	_, err := c.grpc.CloseOracle(ctx, &emptypb.Empty{})
	return err
}

func (c *client) Start(ctx context.Context) error {
	_, err := c.grpc.StartOracle(ctx, &emptypb.Empty{})
	return err
}
