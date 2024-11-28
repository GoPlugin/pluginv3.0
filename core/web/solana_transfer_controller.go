package web

import (
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/goplugin/plugin-common/pkg/types"
	"github.com/goplugin/pluginv3.0/v2/core/chains"
	"github.com/goplugin/pluginv3.0/v2/core/logger/audit"
	"github.com/goplugin/pluginv3.0/v2/core/services/plugin"
	"github.com/goplugin/pluginv3.0/v2/core/services/relay"
	solanamodels "github.com/goplugin/pluginv3.0/v2/core/store/models/solana"
	"github.com/goplugin/pluginv3.0/v2/core/web/presenters"
)

// SolanaTransfersController can send PLI tokens to another address
type SolanaTransfersController struct {
	App plugin.Application
}

// Create sends SOL and other native coins from the Plugin's account to a specified address.
func (tc *SolanaTransfersController) Create(c *gin.Context) {
	relayers := tc.App.GetRelayers().List(plugin.FilterRelayersByType(relay.NetworkSolana))
	if relayers == nil {
		jsonAPIError(c, http.StatusBadRequest, ErrSolanaNotEnabled)
		return
	}

	var tr solanamodels.SendRequest
	if err := c.ShouldBindJSON(&tr); err != nil {
		jsonAPIError(c, http.StatusBadRequest, err)
		return
	}
	if tr.SolanaChainID == "" {
		jsonAPIError(c, http.StatusBadRequest, errors.New("missing solanaChainID"))
		return
	}
	if tr.From.IsZero() {
		jsonAPIError(c, http.StatusUnprocessableEntity, errors.Errorf("source address is missing: %v", tr.From))
		return
	}
	if tr.Amount == 0 {
		jsonAPIError(c, http.StatusBadRequest, errors.New("amount must be greater than zero"))
		return
	}

	amount := new(big.Int).SetUint64(tr.Amount)
	relayerID := types.RelayID{Network: relay.NetworkSolana, ChainID: tr.SolanaChainID}
	relayer, err := relayers.Get(relayerID)
	if err != nil {
		if errors.Is(err, plugin.ErrNoSuchRelayer) {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	err = relayer.Transact(c.Request.Context(), tr.From.String(), tr.To.String(), amount, !tr.AllowHigherAmounts)
	if err != nil {
		if errors.Is(err, chains.ErrNotFound) || errors.Is(err, chains.ErrChainIDEmpty) {
			jsonAPIError(c, http.StatusBadRequest, err)
			return
		}
		jsonAPIError(c, http.StatusInternalServerError, err)
		return
	}

	resource := presenters.NewSolanaMsgResource("sol_transfer_"+uuid.New().String(), tr.SolanaChainID)
	resource.Amount = tr.Amount
	resource.From = tr.From.String()
	resource.To = tr.To.String()

	tc.App.GetAuditLogger().Audit(audit.SolanaTransactionCreated, map[string]interface{}{
		"solanaTransactionResource": resource,
	})
	jsonAPIResponse(c, resource, "solana_tx")
}
