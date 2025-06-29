package handler

import (
	"net/http"
	"root-app/internal/api/dto"
	"root-app/internal/contract"
	"root-app/internal/entities"
	"root-app/internal/exception"
	"root-app/internal/utils"
	"root-app/pkg/gin_helper"
	"root-app/pkg/web_response"

	"github.com/gin-gonic/gin"
)

// EventHandler handles HTTP requests related to event management.
type EventHandler struct {
	eventService contract.EventService
	authService  contract.UserService
}

// NewEventHandler creates and returns a new instance of EventHandler.
func NewEventHandler(eventService contract.EventService, authService contract.UserService) *EventHandler {
	return &EventHandler{
		eventService: eventService,
		authService:  authService,
	}
}

// CreateEvent handles the creation of a new event.
// @Router /api/events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req dto.CreateEventRequest
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
		UserID:    userID.String(),
		Username:  username,
		IPAddress: c.ClientIP(),
	}

	PriceGd, appErr := utils.StringToGormDecimal(req.Price)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	event := &entities.Event{
		Name:        req.Name,
		Description: req.Description,
		Venue:       req.Venue,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Capacity:    req.Capacity,
		// FIX: The entity field `Price` is a pointer. `PriceGd` is already a pointer. Do not dereference it.
		Price: PriceGd,
		// FIX: The entity field `Status` has type `entities.EventStatus`, not `string`. Cast it.
		Status: entities.EventStatus(req.Status),
	}

	createdEvent, appErr := h.eventService.CreateEvent(c.Request.Context(), event, logCtx)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	resp := dto.MapEventToResponse(createdEvent)
	web_response.RespondWithSuccess(c, http.StatusCreated, resp)
}

// GetEventByID handles retrieving a single event by ID.
// @Router /api/events/{id} [get]
func (h *EventHandler) GetEventByID(c *gin.Context) {
	eventID, appErr := gin_helper.ParseIDFromContext(c, "id", "event")
	if appErr != nil {
		return // Response already handled by ParseIDFromContext
	}

	event, appErr := h.eventService.GetEventByID(c.Request.Context(), eventID)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	resp := dto.MapEventToResponse(event)
	web_response.RespondWithSuccess(c, http.StatusOK, resp)
}

// UpdateEvent handles updating an existing event.
// @Router /api/events/{id} [patch]
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	eventID, appErr := gin_helper.ParseIDFromContext(c, "id", "event")
	if appErr != nil {
		return
	}

	var req dto.UpdateEventRequest
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
		UserID:    userID.String(),
		Username:  username,
		IPAddress: c.ClientIP(),
	}

	// FIX: The service expects an *entities.Event object, not an ID.
	// We must construct this object from the request DTO and the ID from the path.
	eventToUpdate := &entities.Event{
		ID: eventID,
	}
	// Populate fields from the request DTO if they are not nil
	if req.Name != nil {
		eventToUpdate.Name = *req.Name
	}
	if req.Description != nil {
		eventToUpdate.Description = *req.Description
	}
	if req.Venue != nil {
		eventToUpdate.Venue = *req.Venue
	}
	if req.StartDate != nil {
		eventToUpdate.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		eventToUpdate.EndDate = *req.EndDate
	}
	if req.Capacity != nil {
		eventToUpdate.Capacity = *req.Capacity
	}
	if req.Status != nil {
		eventToUpdate.Status = entities.EventStatus(*req.Status)
	}
	// Note: Price update logic might be needed here as well.

	updatedEvent, appErr := h.eventService.UpdateEvent(c.Request.Context(), eventToUpdate, logCtx)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	resp := dto.MapEventToResponse(updatedEvent)
	web_response.RespondWithSuccess(c, http.StatusOK, resp)
}

// DeleteEvent handles deleting an event by ID.
// @Router /api/events/{id} [delete]
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	eventID, appErr := gin_helper.ParseIDFromContext(c, "id", "event")
	if appErr != nil {
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
		UserID:    userID.String(),
		Username:  username,
		IPAddress: c.ClientIP(),
	}

	appErr = h.eventService.DeleteEvent(c.Request.Context(), eventID, logCtx)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	web_response.RespondWithSuccess(c, http.StatusNoContent, nil)
}

// GetAllEvents handles listing all events.
// @Router /api/events/ [get]
func (h *EventHandler) GetAllEvents(c *gin.Context) {
	var filter contract.EventFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		// FIX: `web_response.Error` is undefined. Use `HandleAppError`.
		web_response.HandleAppError(c, exception.NewValidationError("Invalid query parameters", err.Error()))
		return
	}

	// FIX: `logCtx` must be created here. It cannot be retrieved from `c.MustGet` unless a middleware sets it.
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

	// FIX: Typo in service method name. It's `GetAllEvents`, not `GetAllTickets`.
	events, total, appErr := h.eventService.GetAllEvents(c.Request.Context(), filter, logCtx)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}

	// FIX: Map entities to DTOs for a clean response structure.
	eventResponses := make([]*dto.EventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = dto.MapEventToResponse(event)
	}

	// FIX: Corrected response format to be consistent and avoid too many arguments.
	web_response.RespondWithSuccess(c, http.StatusOK, gin.H{
		"data": gin.H{
			"events": eventResponses,
			"total":  total,
			"page":   filter.Page,
			"limit":  filter.Limit,
		},
		"message": "Events retrieved successfully",
	})
}