package loadvrfv2plus

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/goplugin/plugin-testing-framework/wasp"

	"github.com/goplugin/pluginv3.0/integration-tests/contracts"
)

/* Monitors on-chain stats of LoadConsumer and pushes them to Loki every second */

const (
	LokiTypeLabel = "vrfv2plus_contracts_load_summary"
	ErrMetrics    = "failed to get VRFv2Plus load test metrics"
	ErrLokiClient = "failed to create Loki client for monitoring"
	ErrLokiPush   = "failed to push monitoring metrics to Loki"
)

func MonitorLoadStats(ctx context.Context, lc *wasp.LokiClient, consumer contracts.VRFv2PlusLoadTestConsumer, labels map[string]string) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			metrics := GetLoadTestMetrics(ctx, consumer)
			SendMetricsToLoki(metrics, lc, labels)
		}
	}()
}

func UpdateLabels(labels map[string]string, t *testing.T) map[string]string {
	updatedLabels := make(map[string]string)
	for k, v := range labels {
		updatedLabels[k] = v
	}
	updatedLabels["type"] = LokiTypeLabel
	updatedLabels["go_test_name"] = t.Name()
	updatedLabels["gen_name"] = "performance"
	return updatedLabels
}

func SendMetricsToLoki(metrics *contracts.VRFLoadTestMetrics, lc *wasp.LokiClient, updatedLabels map[string]string) {
	if err := lc.HandleStruct(wasp.LabelsMapToModel(updatedLabels), time.Now(), metrics); err != nil {
		log.Error().Err(err).Msg(ErrLokiPush)
	}
}

func GetLoadTestMetrics(ctx context.Context, consumer contracts.VRFv2PlusLoadTestConsumer) *contracts.VRFLoadTestMetrics {
	metrics, err := consumer.GetLoadTestMetrics(ctx)
	if err != nil {
		log.Error().Err(err).Msg(ErrMetrics)
	}
	return metrics
}
