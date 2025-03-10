package web

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"regexp"

	"github.com/goplugin/pluginv3.0/v2/core/auth"
	"github.com/goplugin/pluginv3.0/v2/core/bridges"
	"github.com/goplugin/pluginv3.0/v2/core/logger/audit"
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/store/models"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	externalInitiatorNameRegexp = regexp.MustCompile("^[a-zA-Z0-9-_]+$")
)

// ValidateExternalInitiator checks whether External Initiator parameters are
// safe for processing.
func ValidateExternalInitiator(
	ctx context.Context,
	exi *bridges.ExternalInitiatorRequest,
	orm bridges.ORM,
) error {
	fe := models.NewJSONAPIErrors()
	if len([]rune(exi.Name)) == 0 {
		fe.Add("No name specified")
	} else if !externalInitiatorNameRegexp.MatchString(exi.Name) {
		fe.Add("Name must be alphanumeric and may contain '_' or '-'")
	} else if _, err := orm.FindExternalInitiatorByName(ctx, exi.Name); err == nil {
		fe.Add(fmt.Sprintf("Name %v already exists", exi.Name))
	} else if !errors.Is(err, sql.ErrNoRows) {
		return errors.Wrap(err, "validating external initiator")
	}
	return fe.CoerceEmptyToNil()
}

// ExternalInitiatorsController manages external initiators
type ExternalInitiatorsController struct {
	App plugin.Application
}

func (eic *ExternalInitiatorsController) Index(c *gin.Context, size, page, offset int) {
	ctx := c.Request.Context()
	eis, count, err := eic.App.BridgeORM().ExternalInitiators(ctx, offset, size)
	var resources []presenters.ExternalInitiatorResource
	for _, ei := range eis {
		resources = append(resources, presenters.NewExternalInitiatorResource(ei))
	}

	paginatedResponse(c, "externalInitiators", size, page, resources, count, err)
}

// Create builds and saves a new external initiator
func (eic *ExternalInitiatorsController) Create(c *gin.Context) {
	ctx := c.Request.Context()
	eir := &bridges.ExternalInitiatorRequest{}
	if !eic.App.GetConfig().JobPipeline().ExternalInitiatorsEnabled() {
		err := errors.New("The External Initiator feature is disabled by configuration")
		jsonAPIError(c, http.StatusMethodNotAllowed, err)
		return
	}

	eia := auth.NewToken()
	if err := c.ShouldBindJSON(eir); err != nil {
		jsonAPIError(c, http.StatusUnprocessableEntity, err)
		return
	}

	ei, err := bridges.NewExternalInitiator(eia, eir)
	if err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	if err := ValidateExternalInitiator(ctx, eir, eic.App.BridgeORM()); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	if err := eic.App.BridgeORM().CreateExternalInitiator(ctx, ei); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	eic.App.GetAuditLogger().Audit(audit.ExternalInitiatorCreated, map[string]interface{}{
		"externalInitiatorID":   ei.ID,
		"externalInitiatorName": ei.Name,
		"externalInitiatorURL":  ei.URL,
	})

	resp := presenters.NewExternalInitiatorAuthentication(*ei, *eia)
	jsonAPIResponseWithStatus(c, resp, "external initiator authentication", http.StatusCreated)
}

// Destroy deletes an ExternalInitiator
func (eic *ExternalInitiatorsController) Destroy(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("Name")
	exi, err := eic.App.BridgeORM().FindExternalInitiatorByName(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		jsonAPIError(c, http.StatusNotFound, errors.New("external initiator not found"))
		return
	}
	if err := eic.App.BridgeORM().DeleteExternalInitiator(ctx, exi.Name); err != nil {
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	eic.App.GetAuditLogger().Audit(audit.ExternalInitiatorDeleted, map[string]interface{}{"name": name})
	jsonAPIResponseWithStatus(c, nil, "external initiator", http.StatusNoContent)
}
