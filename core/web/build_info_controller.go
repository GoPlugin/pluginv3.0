package web

import (
	"net/http"

	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/static"

	"github.com/gin-gonic/gin"
)

// BuildVersonController has the build_info endpoint.
type BuildInfoController struct {
	App plugin.Application
}

// Show returns the build info.
func (eic *BuildInfoController) Show(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": static.Version, "commitSHA": static.Sha})
}
