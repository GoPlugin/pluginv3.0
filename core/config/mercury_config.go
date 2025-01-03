package config

import (
	"time"

	commonconfig "github.com/goplugin/plugin-common/pkg/config"
	"github.com/goplugin/plugin-common/pkg/types"
)

type MercuryCache interface {
	LatestReportTTL() time.Duration
	MaxStaleAge() time.Duration
	LatestReportDeadline() time.Duration
}

type MercuryTLS interface {
	CertFile() string
}

type MercuryTransmitter interface {
	TransmitQueueMaxSize() uint32
	TransmitTimeout() commonconfig.Duration
}

type Mercury interface {
	Credentials(credName string) *types.MercuryCredentials
	Cache() MercuryCache
	TLS() MercuryTLS
	Transmitter() MercuryTransmitter
	VerboseLogging() bool
}
