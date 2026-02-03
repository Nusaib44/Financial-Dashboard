package handlers

import (
	"net/http"

	"github.com/agency-finance-reality/server/internal/services"
	"github.com/gin-gonic/gin"
)

type UtilizationHandler struct {
	agencyService      services.AgencyService
	utilizationService services.UtilizationService
}

func NewUtilizationHandler(agencyService services.AgencyService, utilizationService services.UtilizationService) *UtilizationHandler {
	return &UtilizationHandler{
		agencyService:      agencyService,
		utilizationService: utilizationService,
	}
}

type AddTimeEntryRequest struct {
	ClientID *string `json:"client_id"`
	Hours    float64 `json:"hours" binding:"required,gt=0"`
}

func (h *UtilizationHandler) AddTimeEntry(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}

	var req AddTimeEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err = h.utilizationService.AddTimeEntry(agency.ID, req.ClientID, req.Hours)
	if err != nil {
		SendInternalError(c)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *UtilizationHandler) GetUtilization(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}

	utilization, err := h.utilizationService.GetUtilization(agency.ID)
	if err != nil {
		SendInternalError(c)
		return
	}

	c.JSON(http.StatusOK, utilization)
}
