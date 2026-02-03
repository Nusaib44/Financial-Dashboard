package handlers

import (
	"net/http"

	"github.com/agency-finance-reality/server/internal/services"
	"github.com/gin-gonic/gin"
)

type SurvivalHandler struct {
	agencyService  services.AgencyService
	financeService services.FinanceService
}

func NewSurvivalHandler(agencyService services.AgencyService, financeService services.FinanceService) *SurvivalHandler {
	return &SurvivalHandler{
		agencyService:  agencyService,
		financeService: financeService,
	}
}

func (h *SurvivalHandler) GetBurnRunway(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}

	metrics, err := h.financeService.GetSurvivalMetrics(agency.ID)
	if err != nil {
		SendInternalError(c)
		return
	}

	c.JSON(http.StatusOK, metrics)
}
