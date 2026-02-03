package handlers

import (
	"net/http"

	"github.com/agency-finance-reality/server/internal/services"
	"github.com/gin-gonic/gin"
)

type RealityScoreHandler struct {
	agencyService  services.AgencyService
	financeService services.FinanceService
}

func NewRealityScoreHandler(agencyService services.AgencyService, financeService services.FinanceService) *RealityScoreHandler {
	return &RealityScoreHandler{
		agencyService:  agencyService,
		financeService: financeService,
	}
}

func (h *RealityScoreHandler) GetRealityScore(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}

	score, err := h.financeService.GetRealityScore(agency.ID)
	if err != nil {
		SendInternalError(c)
		return
	}

	c.JSON(http.StatusOK, score)
}
