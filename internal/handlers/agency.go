package handlers

import (
	"net/http"

	"github.com/agency-finance-reality/server/internal/services"
	"github.com/gin-gonic/gin"
)

type AgencyHandler struct {
	agencyService services.AgencyService
}

func NewAgencyHandler(agencyService services.AgencyService) *AgencyHandler {
	return &AgencyHandler{agencyService: agencyService}
}

func (h *AgencyHandler) GetAgency(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}
	c.JSON(http.StatusOK, agency)
}

type CreateAgencyRequest struct {
	Name         string  `json:"name" binding:"required"`
	BaseCurrency string  `json:"base_currency" binding:"required"`
	StartingCash float64 `json:"starting_cash"`
}

func (h *AgencyHandler) CreateAgency(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	var req CreateAgencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err := h.agencyService.CreateAgency(userID, req.Name, req.BaseCurrency, req.StartingCash)
	if err != nil {
		SendInternalError(c)
		return
	}

	c.Status(http.StatusCreated)
}
