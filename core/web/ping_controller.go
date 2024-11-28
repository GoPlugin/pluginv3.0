package web

import (
	"net/http"

	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"

	"github.com/gin-gonic/gin"
)

// PingController has the ping endpoint.
type PingController struct {
	App plugin.Application
}

// Show returns pong.
func (eic *PingController) Show(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
