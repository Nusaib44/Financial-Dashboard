package handlers

import (
	"net/http"

	"github.com/agency-finance-reality/server/internal/services"
	"github.com/gin-gonic/gin"
)

type CashSnapshotHandler struct {
	agencyService  services.AgencyService
	financeService services.FinanceService
}

func NewCashSnapshotHandler(agencyService services.AgencyService, financeService services.FinanceService) *CashSnapshotHandler {
	return &CashSnapshotHandler{
		agencyService:  agencyService,
		financeService: financeService,
	}
}

func (h *CashSnapshotHandler) GetTodaysCash(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}

	snap, err := h.financeService.GetDailySnapshot(agency.ID)
	if err != nil {
		SendInternalError(c)
		return
	}
	if snap == nil {
		SendError(c, http.StatusNotFound, "No snapshot today")
		return
	}

	c.JSON(http.StatusOK, snap)
}

type RecordCashRequest struct {
	CashBalance float64 `json:"cash_balance" binding:"required"`
}

func (h *CashSnapshotHandler) RecordDailyCash(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}

	var req RecordCashRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	err = h.financeService.RecordCashSnapshot(agency.ID, req.CashBalance)
	if err != nil {
		SendInternalError(c)
		return
	}

	c.Status(http.StatusCreated)
}
