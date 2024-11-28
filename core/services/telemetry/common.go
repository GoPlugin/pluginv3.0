package telemetry

import (
	ocrtypes "github.com/goplugin/plugin-libocr/commontypes"

	"github.com/goplugin/pluginv3.0/v2/core/services/synchronization"
)

type MonitoringEndpointGenerator interface {
	GenMonitoringEndpoint(network string, chainID string, contractID string, telemType synchronization.TelemetryType) ocrtypes.MonitoringEndpoint
}
