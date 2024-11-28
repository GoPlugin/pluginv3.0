package web

import (
	"github.com/gin-gonic/gin"

	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

// FeaturesController manages the feature flags
type FeaturesController struct {
	App plugin.Application
}

const (
	FeatureKeyCSA                string = "csa"
	FeatureKeyFeedsManager       string = "feeds_manager"
	FeatureKeyMultiFeedsManagers string = "multi_feeds_managers"
)

// Index retrieves the features
// Example:
// "GET <application>/features"
func (fc *FeaturesController) Index(c *gin.Context) {
	resources := []presenters.FeatureResource{
		*presenters.NewFeatureResource(FeatureKeyCSA, fc.App.GetConfig().Feature().UICSAKeys()),
		*presenters.NewFeatureResource(FeatureKeyFeedsManager, fc.App.GetConfig().Feature().FeedsManager()),
		*presenters.NewFeatureResource(FeatureKeyMultiFeedsManagers, fc.App.GetConfig().Feature().MultiFeedsManagers()),
	}

	jsonAPIResponse(c, resources, "features")
}
