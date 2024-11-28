package telemetry

import (
	ocrtypes "github.com/goplugin/plugin-libocr/commontypes"

	"github.com/goplugin/pluginv3.0/v2/core/services/synchronization"
)

var _ MonitoringEndpointGenerator = &NoopAgent{}

type NoopAgent struct {
}

// SendLog sends a telemetry log to the ingress service
func (t *NoopAgent) SendLog(log []byte) {
}

// GenMonitoringEndpoint creates a monitoring endpoint for telemetry
func (t *NoopAgent) GenMonitoringEndpoint(network string, chainID string, contractID string, telemType synchronization.TelemetryType) ocrtypes.MonitoringEndpoint {
	return t
}
