package managed

import (
	"context"

	"github.com/goplugin/plugin-libocr/commontypes"
	"github.com/goplugin/plugin-libocr/internal/loghelper"
	"github.com/goplugin/plugin-libocr/offchainreporting/internal/serialization/protobuf"
	"google.golang.org/protobuf/proto"
)

// forwardTelemetry receives monitoring events on chTelemetry, serializes them, and forwards
// them to monitoringEndpoint
func forwardTelemetry(
	ctx context.Context,

	logger loghelper.LoggerWithContext,
	monitoringEndpoint commontypes.MonitoringEndpoint,

	chTelemetry <-chan *protobuf.TelemetryWrapper,
) {
	for {
		select {
		case t, ok := <-chTelemetry:
			if !ok {
				// This isn't supposed to happen, but we still handle this case gracefully,
				// just in case...
				logger.Error("forwardTelemetry: chTelemetry closed unexpectedly. exiting", nil)
				return
			}
			bin, err := proto.Marshal(t)
			if err != nil {
				logger.Error("forwardTelemetry: failed to Marshal protobuf", commontypes.LogFields{
					"proto": t,
					"error": err,
				})
				break
			}
			if monitoringEndpoint != nil {
				monitoringEndpoint.SendLog(bin)
			}
		case <-ctx.Done():
			logger.Info("forwardTelemetry: exiting", nil)
			return
		}
	}
}
