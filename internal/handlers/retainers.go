package handlers

import (
	"net/http"

	"github.com/agency-finance-reality/server/internal/services"
	"github.com/gin-gonic/gin"
)

type RetainerHandler struct {
	agencyService services.AgencyService
	clientService services.ClientService
}

func NewRetainerHandler(agencyService services.AgencyService, clientService services.ClientService) *RetainerHandler {
	return &RetainerHandler{
		agencyService: agencyService,
		clientService: clientService,
	}
}

type CreateRetainerRequest struct {
	ClientID      string  `json:"client_id" binding:"required"`
	MonthlyAmount float64 `json:"monthly_amount" binding:"required,gt=0"`
}

func (h *RetainerHandler) CreateRetainer(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}

	var req CreateRetainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err = h.clientService.CreateRetainer(agency.ID, req.ClientID, req.MonthlyAmount)
	if err != nil {
		SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *RetainerHandler) GetRetainerSummary(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}

	summary, err := h.clientService.GetRetainerSummary(agency.ID)
	if err != nil {
		SendInternalError(c)
		return
	}

	c.JSON(http.StatusOK, summary)
}
