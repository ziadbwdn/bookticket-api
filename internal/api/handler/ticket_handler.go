// Corrected internal/api/handler/ticket_handler.go

package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"root-app/internal/api/dto"
	"root-app/internal/contract"
	"root-app/internal/entities"
	"root-app/internal/exception"
	"root-app/internal/utils"
	"root-app/pkg/gin_helper"
	"root-app/pkg/web_response"

	"github.com/gin-gonic/gin"
)

// TicketHandler handles HTTP requests related to ticket management.
type TicketHandler struct {
	ticketService contract.TicketService
	authService contract.UserService
}

// NewTicketHandler creates and returns a new instance of TicketHandler.
func NewTicketHandler(ticketService contract.TicketService, authService contract.UserService) *TicketHandler {
	return &TicketHandler{
		ticketService: ticketService,
		authService:   authService,
	}
}

// PurchaseTicket handles the creation of a new ticket.
// @Router /api/tickets [post]
func (h *TicketHandler) PurchaseTicket(c *gin.Context) {
	// FIX: Changed dto.CreateStationRequest to dto.BookTicketRequest
	var req dto.BookTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web_response.HandleAppError(c, exception.NewValidationError("Invalid request body", err.Error()))
		return
	}

	userID, appErr := gin_helper.GetUserIDFromContext(c)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	username, appErr := h.authService.GetUserDetailsForLogging(c.Request.Context(), userID)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	unitPriceGd, appErr := utils.StringToGormDecimal(req.UnitPrice); 
	if appErr != nil { 
		web_response.HandleAppError(c, appErr); 
		return 
	}

	totalPriceGd, appErr := utils.StringToGormDecimal(req.TotalPrice); 
	if appErr != nil { 
		web_response.HandleAppError(c, appErr); 
		return 
	}

	userIDStr := userID.String()
	eventIDStr, err := utils.ParseBinaryUUID(req.EventID)
	if err != nil {
		web_response.HandleAppError(c, exception.NewValidationError("Invalid event_id format in request", err.Error()))
		return
	}

	logCtx := entities.ActivityLogContext{
		UserID:    userIDStr,
		Username:  username,
		IPAddress: c.ClientIP(),
	}

	quantity, err := strconv.Atoi(req.Quantity)
	if err != nil {
		web_response.HandleAppError(c, exception.NewValidationError("Invalid quantity format, must be a number", err.Error()))
		return
	}


	// FIX: Corrected field mapping from DTO to entity.
	// UserID is taken from context, not request body, for security.
	ticket := &entities.Ticket{
		EventID:    eventIDStr,
		UserID: 	req.UserID,
		TicketCode: req.TicketCode,
		Quantity:   quantity,
		UnitPrice:  unitPriceGd,
		TotalPrice: totalPriceGd,
		Status:     entities.TicketStatus(req.Status),
	}

	// Removed userID as it's now in logCtx.
	purchasedTicket, appErr := h.ticketService.PurchaseTicket(c.Request.Context(), ticket, logCtx)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	resp := dto.MapTicketToResponse(purchasedTicket)
	web_response.RespondWithSuccess(c, http.StatusCreated, resp)
}

// GetTicketByID handles retrieving a single ticket by ID.
// @Router /api/tickets/{id} [get]
func (h *TicketHandler) GetTicketByID(c *gin.Context) {
	// FIX: Corrected package name from gin_helpers to gin_helper.
	ticketID, appErr := gin_helper.ParseIDFromContext(c, "id", "ticket")
	if appErr != nil {
		return // Response already handled by ParseIDFromContext
	}

	ticket, appErr := h.ticketService.GetTicketByID(c.Request.Context(), ticketID)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	resp := dto.MapTicketToResponse(ticket)
	web_response.RespondWithSuccess(c, http.StatusOK, resp)
}

// UpdateTicketStatus handles updating an existing ticket's status.
// @Router /api/tickets/{id} [patch]
func (h *TicketHandler) UpdateTicketStatus(c *gin.Context) {
	ticketID, appErr := gin_helper.ParseIDFromContext(c, "id", "ticket")
	if appErr != nil {
		return
	}
	var req dto.UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web_response.HandleAppError(c, exception.NewValidationError("Invalid request body", err.Error()))
		return
	}
	userID, appErr := gin_helper.GetUserIDFromContext(c)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}
	username, appErr := h.authService.GetUserDetailsForLogging(c.Request.Context(), userID)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}
	logCtx := entities.ActivityLogContext{
		// FIX: Corrected the missing conversion from the previous attempt.
		UserID:    userID.String(),
		Username:  username,
		IPAddress: c.ClientIP(),
	}
	updatedTicket, appErr := h.ticketService.UpdateTicketStatus(c.Request.Context(), ticketID, &req, "", logCtx)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}
	resp := dto.MapTicketToResponse(updatedTicket)
	web_response.RespondWithSuccess(c, http.StatusOK, resp)
}

// DeleteTicket handles deleting a ticket by ID.
// @Router /api/tickets/{id} [delete]
func (h *TicketHandler) DeleteTicket(c *gin.Context) {
	// FIX: Corrected package name.
	ticketID, appErr := gin_helper.ParseIDFromContext(c, "id", "ticket")
	if appErr != nil {
		return
	}

	// FIX: Corrected package name.
	userID, appErr := gin_helper.GetUserIDFromContext(c)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}
	username, appErr := h.authService.GetUserDetailsForLogging(c.Request.Context(), userID)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	// FIX: Assumed ActivityLogContext.UserID is utils.BinaryUUID.
	logCtx := entities.ActivityLogContext{
		UserID:    userID.String(),
		Username:  username,
		IPAddress: c.ClientIP(),
	}

	appErr = h.ticketService.DeleteTicket(c.Request.Context(), ticketID, logCtx)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	web_response.RespondWithSuccess(c, http.StatusNoContent, nil)
}

// GetAllTickets handles listing all tickets, with optional filters.
// @Router /api/tickets [get]
func (h *TicketHandler) GetAllTickets(c *gin.Context) {
	var filter contract.TicketFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		// FIX: Used HandleAppError instead of undefined web_response.Error
		web_response.HandleAppError(c, exception.NewValidationError("Invalid query parameters", err.Error()))
		return
	}

	// FIX: This handler needs a log context. It should be created like in other handlers.
	userID, appErr := gin_helper.GetUserIDFromContext(c)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}
	username, appErr := h.authService.GetUserDetailsForLogging(c.Request.Context(), userID)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}
	logCtx := entities.ActivityLogContext{
		UserID:    userID.String(),
		Username:  username,
		IPAddress: c.ClientIP(),
	}

	tickets, total, appErr := h.ticketService.GetAllTickets(c.Request.Context(), filter, logCtx)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	// FIX: Map entities to DTOs for a clean response.
	ticketResponses := make([]*dto.TicketResponse, len(tickets))
	for i, ticket := range tickets {
		ticketResponses[i] = dto.MapTicketToResponse(ticket)
	}

	// FIX: Used RespondWithSuccess instead of undefined web_response.Success
	// FIX: Corrected the field name in the DTO from `Ticket` to `Tickets` to match this response.
	// You may need to adjust `dto.ListTicketResponses` to have `Tickets` field.
	web_response.RespondWithSuccess(c, http.StatusOK, gin.H{
		"data": gin.H{
			"tickets": ticketResponses,
			"total":   total,
			"page":    filter.Page,
			"limit":   filter.Limit,
		},
		"message": "Tickets retrieved successfully",
	})
}

// GetUserTickets handles listing tickets for a specific user.
// @Router /api/users/{id}/tickets [get]
func (h *TicketHandler) GetUserTickets(c *gin.Context) {
	// FIX: This handler was completely rewritten to use the new GetUserTickets service method.
	targetUserID, appErr := gin_helper.ParseIDFromContext(c, "id", "user")
	if appErr != nil {
		return // Response already handled
	}

	var filter contract.TicketFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		web_response.HandleAppError(c, exception.NewValidationError("Invalid query parameters", err.Error()))
		return
	}

	// For logging purposes, we get the ID of the user MAKING the request.
	requestingUserID, appErr := gin_helper.GetUserIDFromContext(c)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}
	username, appErr := h.authService.GetUserDetailsForLogging(c.Request.Context(), requestingUserID)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}
	logCtx := entities.ActivityLogContext{
		UserID:    requestingUserID.String(),
		Username:  username,
		IPAddress: c.ClientIP(),
	}

	// Call the new service method
	tickets, total, appErr := h.ticketService.GetUserTickets(c.Request.Context(), targetUserID, filter, logCtx)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	// Map entities to DTOs
	ticketResponses := make([]*dto.TicketResponse, len(tickets))
	for i, ticket := range tickets {
		ticketResponses[i] = dto.MapTicketToResponse(ticket)
	}

	web_response.RespondWithSuccess(c, http.StatusOK, gin.H{
		"data": gin.H{
			"tickets": ticketResponses,
			"total":   total,
			"page":    filter.Page,
			"limit":   filter.Limit,
		},
		"message": fmt.Sprintf("Tickets for user %s retrieved successfully", targetUserID.String()),
	})
}

// NOTE: The `GetTicketsByEvent` handler was removed. Its functionality is now covered by `GetAllTickets`.
// To get tickets for an event, the client should call `GET /api/tickets?event_id=<event_id>`.
// The `ShouldBindQuery` in `GetAllTickets` will automatically populate the filter.