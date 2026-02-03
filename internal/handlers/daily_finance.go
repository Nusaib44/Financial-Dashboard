package handlers

import (
	"net/http"

	"github.com/agency-finance-reality/server/internal/services"
	"github.com/gin-gonic/gin"
)

type DailyFinanceHandler struct {
	agencyService  services.AgencyService
	financeService services.FinanceService
}

func NewDailyFinanceHandler(agencyService services.AgencyService, financeService services.FinanceService) *DailyFinanceHandler {
	return &DailyFinanceHandler{
		agencyService:  agencyService,
		financeService: financeService,
	}
}

type AddRevenueRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Source string  `json:"source" binding:"required"`
}

func (h *DailyFinanceHandler) AddRevenue(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}

	var req AddRevenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err = h.financeService.AddRevenue(agency.ID, req.Amount, req.Source)
	if err != nil {
		SendInternalError(c)
		return
	}

	c.Status(http.StatusCreated)
}

type AddCostRequest struct {
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Type     string  `json:"type" binding:"required"`
	Label    string  `json:"label" binding:"required"`
	Category string  `json:"category" binding:"required,oneof=people tools other"`
}

func (h *DailyFinanceHandler) AddCost(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}

	var req AddCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err = h.financeService.AddCost(agency.ID, req.Amount, req.Type, req.Label, req.Category)
	if err != nil {
		SendInternalError(c)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *DailyFinanceHandler) GetCostBreakdown(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}

	breakdown, err := h.financeService.GetCostBreakdown(agency.ID)
	if err != nil {
		SendInternalError(c)
		return
	}

	c.JSON(http.StatusOK, breakdown)
}

func (h *DailyFinanceHandler) GetDailySummary(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}

	summary, err := h.financeService.GetDailySummary(agency.ID)
	if err != nil {
		SendInternalError(c)
		return
	}

	c.JSON(http.StatusOK, summary)
}
