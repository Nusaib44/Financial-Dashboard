package handlers

import (
	"net/http"

	"github.com/agency-finance-reality/server/internal/services"
	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	agencyService services.AgencyService
	clientService services.ClientService
}

func NewClientHandler(agencyService services.AgencyService, clientService services.ClientService) *ClientHandler {
	return &ClientHandler{
		agencyService: agencyService,
		clientService: clientService,
	}
}

type CreateClientRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *ClientHandler) CreateClient(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}

	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	client, err := h.clientService.CreateClient(agency.ID, req.Name)
	if err != nil {
		SendInternalError(c)
		return
	}

	c.JSON(http.StatusCreated, client)
}

func (h *ClientHandler) GetClients(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	agency, err := h.agencyService.GetAgencyByUserID(userID)
	if err != nil {
		SendError(c, http.StatusNotFound, "Agency not found")
		return
	}

	clients, err := h.clientService.GetClients(agency.ID)
	if err != nil {
		SendInternalError(c)
		return
	}

	c.JSON(http.StatusOK, clients)
}
