package v4

import (
	"context"

	ocr2plus_types "github.com/goplugin/plugin-libocr/offchainreporting2plus/types"
	"google.golang.org/grpc"

	"github.com/goplugin/plugin-common/pkg/loop/internal/pb"
	mercury_v4_pb "github.com/goplugin/plugin-common/pkg/loop/internal/pb/mercury/v4"
	mercury_v4_types "github.com/goplugin/plugin-common/pkg/types/mercury/v4"
)

var _ mercury_v4_types.ReportCodec = (*ReportCodecClient)(nil)

type ReportCodecClient struct {
	grpc mercury_v4_pb.ReportCodecClient
}

func NewReportCodecClient(cc grpc.ClientConnInterface) *ReportCodecClient {
	return &ReportCodecClient{grpc: mercury_v4_pb.NewReportCodecClient(cc)}
}

func (r *ReportCodecClient) BuildReport(ctx context.Context, fields mercury_v4_types.ReportFields) (ocr2plus_types.Report, error) {
	reply, err := r.grpc.BuildReport(ctx, &mercury_v4_pb.BuildReportRequest{
		ReportFields: pbReportFields(fields),
	})
	if err != nil {
		return ocr2plus_types.Report{}, err
	}
	return reply.Report, nil
}

func (r *ReportCodecClient) MaxReportLength(ctx context.Context, n int) (int, error) {
	reply, err := r.grpc.MaxReportLength(ctx, &mercury_v4_pb.MaxReportLengthRequest{})
	if err != nil {
		return 0, err
	}
	return int(reply.MaxReportLength), nil
}

func (r *ReportCodecClient) ObservationTimestampFromReport(ctx context.Context, report ocr2plus_types.Report) (uint32, error) {
	reply, err := r.grpc.ObservationTimestampFromReport(ctx, &mercury_v4_pb.ObservationTimestampFromReportRequest{
		Report: report,
	})
	if err != nil {
		return 0, err
	}
	return reply.Timestamp, nil
}

func pbReportFields(fields mercury_v4_types.ReportFields) *mercury_v4_pb.ReportFields {
	return &mercury_v4_pb.ReportFields{
		ValidFromTimestamp: fields.ValidFromTimestamp,
		Timestamp:          fields.Timestamp,
		NativeFee:          pb.NewBigIntFromInt(fields.NativeFee),
		LinkFee:            pb.NewBigIntFromInt(fields.LinkFee),
		ExpiresAt:          fields.ExpiresAt,
		BenchmarkPrice:     pb.NewBigIntFromInt(fields.BenchmarkPrice),
		MarketStatus:       fields.MarketStatus,
	}
}

var _ mercury_v4_pb.ReportCodecServer = (*ReportCodecServer)(nil)

type ReportCodecServer struct {
	mercury_v4_pb.UnimplementedReportCodecServer
	impl mercury_v4_types.ReportCodec
}

func NewReportCodecServer(impl mercury_v4_types.ReportCodec) *ReportCodecServer {
	return &ReportCodecServer{impl: impl}
}

func (r *ReportCodecServer) BuildReport(ctx context.Context, request *mercury_v4_pb.BuildReportRequest) (*mercury_v4_pb.BuildReportReply, error) {
	report, err := r.impl.BuildReport(ctx, reportFields(request.ReportFields))
	if err != nil {
		return nil, err
	}
	return &mercury_v4_pb.BuildReportReply{Report: report}, nil
}

func (r *ReportCodecServer) MaxReportLength(ctx context.Context, request *mercury_v4_pb.MaxReportLengthRequest) (*mercury_v4_pb.MaxReportLengthReply, error) {
	n, err := r.impl.MaxReportLength(ctx, int(request.NumOracles))
	if err != nil {
		return nil, err
	}
	return &mercury_v4_pb.MaxReportLengthReply{MaxReportLength: uint64(n)}, nil
}

func (r *ReportCodecServer) ObservationTimestampFromReport(ctx context.Context, request *mercury_v4_pb.ObservationTimestampFromReportRequest) (*mercury_v4_pb.ObservationTimestampFromReportReply, error) {
	timestamp, err := r.impl.ObservationTimestampFromReport(ctx, request.Report)
	if err != nil {
		return nil, err
	}
	return &mercury_v4_pb.ObservationTimestampFromReportReply{Timestamp: timestamp}, nil
}

func reportFields(fields *mercury_v4_pb.ReportFields) mercury_v4_types.ReportFields {
	return mercury_v4_types.ReportFields{
		ValidFromTimestamp: fields.ValidFromTimestamp,
		Timestamp:          fields.Timestamp,
		NativeFee:          fields.NativeFee.Int(),
		LinkFee:            fields.LinkFee.Int(),
		ExpiresAt:          fields.ExpiresAt,
		BenchmarkPrice:     fields.BenchmarkPrice.Int(),
		MarketStatus:       fields.MarketStatus,
	}
}
