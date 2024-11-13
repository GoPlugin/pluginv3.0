package plugin

import "github.com/goplugin/pluginv3.0/v2/core/config/toml"

type thresholdConfig struct {
	s toml.ThresholdKeyShareSecrets
}

func (t *thresholdConfig) ThresholdKeyShare() string {
	if t.s.ThresholdKeyShare == nil {
		return ""
	}
	return string(*t.s.ThresholdKeyShare)
}
