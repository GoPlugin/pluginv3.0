package synchronization_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/goplugin/pluginv3.0/v2/core/logger"
	"github.com/goplugin/pluginv3.0/v2/core/services/synchronization"
	"github.com/goplugin/pluginv3.0/v2/core/services/synchronization/mocks"
)

func TestTelemetryIngressWorker_BuildTelemBatchReq(t *testing.T) {
	telemPayload := synchronization.TelemPayload{
		Telemetry:  []byte("Mock telemetry"),
		ContractID: "0xa",
	}

	maxTelemBatchSize := 3
	chTelemetry := make(chan synchronization.TelemPayload, 10)
	worker := synchronization.NewTelemetryIngressBatchWorker(
		uint(maxTelemBatchSize),
		time.Second,
		mocks.NewTelemClient(t),
		chTelemetry,
		"0xa",
		synchronization.OCR,
		logger.TestLogger(t),
		false,
		"test-endpoint",
	)

	chTelemetry <- telemPayload
	chTelemetry <- telemPayload
	chTelemetry <- telemPayload
	chTelemetry <- telemPayload
	chTelemetry <- telemPayload

	// Batch request should not exceed the max batch size
	batchReq1 := worker.BuildTelemBatchReq()
	assert.Equal(t, "0xa", batchReq1.ContractId)
	assert.Equal(t, string(synchronization.OCR), batchReq1.TelemetryType)
	assert.Len(t, batchReq1.Telemetry, maxTelemBatchSize)
	assert.Len(t, chTelemetry, 2)
	assert.Greater(t, batchReq1.SentAt, int64(0))

	// Remainder of telemetry should be batched on next call
	batchReq2 := worker.BuildTelemBatchReq()
	assert.Equal(t, "0xa", batchReq2.ContractId)
	assert.Equal(t, string(synchronization.OCR), batchReq2.TelemetryType)
	assert.Len(t, batchReq2.Telemetry, 2)
	assert.Len(t, chTelemetry, 0)
	assert.Greater(t, batchReq2.SentAt, int64(0))
}
