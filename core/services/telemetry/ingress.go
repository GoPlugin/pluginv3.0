package telemetry

import (
	"context"

	ocrtypes "github.com/goplugin/plugin-libocr/commontypes"

	"github.com/goplugin/pluginv3.0/v2/core/services/synchronization"
)

var _ MonitoringEndpointGenerator = &IngressAgentWrapper{}

type IngressAgentWrapper struct {
	telemetryIngressClient synchronization.TelemetryService
}

func NewIngressAgentWrapper(telemetryIngressClient synchronization.TelemetryService) *IngressAgentWrapper {
	return &IngressAgentWrapper{telemetryIngressClient}
}

func (t *IngressAgentWrapper) GenMonitoringEndpoint(network, chainID string, contractID string, telemType synchronization.TelemetryType) ocrtypes.MonitoringEndpoint {
	return NewIngressAgent(t.telemetryIngressClient, network, chainID, contractID, telemType)
}

type IngressAgent struct {
	telemetryIngressClient synchronization.TelemetryService
	network                string
	chainID                string
	contractID             string
	telemType              synchronization.TelemetryType
}

func NewIngressAgent(telemetryIngressClient synchronization.TelemetryService, network string, chainID string, contractID string, telemType synchronization.TelemetryType) *IngressAgent {
	return &IngressAgent{
		telemetryIngressClient,
		network,
		chainID,
		contractID,
		telemType,
	}
}

// SendLog sends a telemetry log to the ingress server
func (t *IngressAgent) SendLog(telemetry []byte) {
	t.telemetryIngressClient.Send(context.Background(), telemetry, t.contractID, t.telemType)
}
