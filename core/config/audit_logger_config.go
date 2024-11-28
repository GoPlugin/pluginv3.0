package config

import (
	commonconfig "github.com/goplugin/plugin-common/pkg/config"
	"github.com/goplugin/pluginv3.0/v2/core/store/models"
)

type AuditLogger interface {
	Enabled() bool
	ForwardToUrl() (commonconfig.URL, error)
	Environment() string
	JsonWrapperKey() string
	Headers() (models.ServiceHeaders, error)
}
