package plugin

import (
	commonconfig "github.com/goplugin/plugin-common/pkg/config"
	"github.com/goplugin/pluginv3.0/v2/core/build"
	"github.com/goplugin/pluginv3.0/v2/core/config/toml"
	"github.com/goplugin/pluginv3.0/v2/core/store/models"
)

type auditLoggerConfig struct {
	c toml.AuditLogger
}

func (a auditLoggerConfig) Enabled() bool {
	return *a.c.Enabled
}

func (a auditLoggerConfig) ForwardToUrl() (commonconfig.URL, error) {
	return *a.c.ForwardToUrl, nil
}

func (a auditLoggerConfig) Environment() string {
	if !build.IsProd() {
		return "develop"
	}
	return "production"
}

func (a auditLoggerConfig) JsonWrapperKey() string {
	return *a.c.JsonWrapperKey
}

func (a auditLoggerConfig) Headers() (models.ServiceHeaders, error) {
	return *a.c.Headers, nil
}
