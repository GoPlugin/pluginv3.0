package plugin

import (
	"github.com/goplugin/pluginv3.0/v2/core/config/toml"
)

type sentryConfig struct {
	c toml.Sentry
}

func (s sentryConfig) DSN() string {
	return *s.c.DSN
}

func (s sentryConfig) Debug() bool {
	return *s.c.Debug
}

func (s sentryConfig) Environment() string {
	return *s.c.Environment
}

func (s sentryConfig) Release() string {
	return *s.c.Release
}
