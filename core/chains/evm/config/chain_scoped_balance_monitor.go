package config

import "github.com/goplugin/pluginv3.0/v2/core/chains/evm/config/toml"

type balanceMonitorConfig struct {
	c toml.BalanceMonitor
}

func (b *balanceMonitorConfig) Enabled() bool {
	return *b.c.Enabled
}
