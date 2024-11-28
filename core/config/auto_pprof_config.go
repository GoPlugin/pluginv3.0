package config

import (
	commonconfig "github.com/goplugin/plugin-common/pkg/config"
	"github.com/goplugin/pluginv3.0/v2/core/utils"
)

type AutoPprof interface {
	BlockProfileRate() int
	CPUProfileRate() int
	Enabled() bool
	GatherDuration() commonconfig.Duration
	GatherTraceDuration() commonconfig.Duration
	GoroutineThreshold() int
	MaxProfileSize() utils.FileSize
	MemProfileRate() int
	MemThreshold() utils.FileSize
	MutexProfileFraction() int
	PollInterval() commonconfig.Duration
	ProfileRoot() string
}
