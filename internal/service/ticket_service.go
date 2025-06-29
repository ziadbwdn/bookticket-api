// Corrected internal/service/ticket_service.go

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"root-app/internal/api/dto"
	"root-app/internal/contract"
	"root-app/internal/entities"
	"root-app/internal/exception"
	"root-app/internal/logger"
	"root-app/internal/utils"
	"time"
)

// Ticket Service Struct implementation, containing of interfaces from repo
type TicketServiceImpl struct {
	ticketRepo      contract.TicketRepository
	eventRepo       contract.EventRepository
	activityService contract.UserActivityService
	logger          logger.Logger
}

// Ticket Service Constructor
func NewTicketService(ticketRepo contract.TicketRepository, eventRepo contract.EventRepository, activityService contract.UserActivityService, logger logger.Logger) contract.TicketService {
	// panic checking
	if ticketRepo == nil {
		panic("ticketRepo must not be nil for TicketService")
	}
	if eventRepo == nil {
		panic("eventRepo must not be nil for TicketService")
	}
	if activityService == nil {
		panic("activityService must not be nil for TicketService")
	}
	if logger == nil {
		panic("logger must not be nil for TicketService")
	}
	return &TicketServiceImpl{
		ticketRepo:      ticketRepo,
		eventRepo:       eventRepo,
		activityService: activityService,
		logger:          logger,
	}
}

// FIX: The signature was changed to match the contract.TicketService interface.
// The `userID` parameter was removed as it's already inside `logCtx`.
func (s *TicketServiceImpl) PurchaseTicket(ctx context.Context, ticket *entities.Ticket, logCtx entities.ActivityLogContext) (*entities.Ticket, *exception.AppError) {
	ticket.ID = utils.NewBinaryUUID()

	// FIX: Parse the string UserID from logCtx back into a BinaryUUID.
	parsedUserID, err := utils.ParseBinaryUUID(logCtx.UserID)
	if err != nil {
		// This indicates a programming error (a non-UUID string was passed in logCtx)
		s.logger.Error(ctx, "Failed to parse UserID from log context", err, logger.Field{Key: "userID", Value: logCtx.UserID})
		return nil, exception.NewInternalError("invalid_user_id_in_context", err)
	}
	ticket.UserID = parsedUserID // Assign the parsed UUID

	if ticket.CreatedAt.IsZero() {
		ticket.CreatedAt = time.Now()
	}
	ticket.UpdatedAt = time.Now()

	ticket.Status = entities.TicketStatusActive

	appErr := s.ticketRepo.Create(ctx, ticket)
	if appErr != nil {
		return nil, appErr
	}

	// --- LOG USER ACTIVITY ---
	newValueJSON, _ := json.Marshal(ticket)
	newValueStr := string(newValueJSON)
	ticketIDStr := ticket.ID.String()

	// FIX: Corrected undefined field. The ticket code is on the ticket object itself.
	details := fmt.Sprintf("New Ticket '%s' purchase issued. Ticket Code: %s", ticket.ID.String(), ticket.TicketCode)
	ipAddr := logCtx.IPAddress

	// The LogUserActivity service expects a string UserID, which logCtx.UserID already is.
	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, entities.ActionTypePurchaseTicket, entities.ResourceTypeTicket, &ticketIDStr, &ipAddr, &details, nil, &newValueStr)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log PurchaseTicket activity", logErr, logger.Field{Key: "ticketID", Value: ticketIDStr})
	}

	return ticket, nil
}

// Get Ticket by ID
func (s *TicketServiceImpl) GetTicketByID(
	ctx context.Context,
	id utils.BinaryUUID,
) (*entities.Ticket, *exception.AppError) {
	ticket, appErr := s.ticketRepo.GetByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}

	return ticket, nil
}

// FIX: The signature was changed to match the contract.TicketService interface.
// It now takes `id`, `req`, `reason`, and `logCtx`.
// The return type was changed from *entities.Event to *entities.Ticket.
func (s *TicketServiceImpl) UpdateTicketStatus(ctx context.Context, id utils.BinaryUUID, req *dto.UpdateTicketStatusRequest, reason string, logCtx entities.ActivityLogContext) (*entities.Ticket, *exception.AppError) {
	// Fetching current record from the db
	existingTicket, appErr := s.ticketRepo.GetByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}

	// Marshal the original state for logging before any changes are made.
	oldValueJSON, err := json.Marshal(existingTicket)
	if err != nil {
		s.logger.Warn(ctx, "Failed to marshal old ticket value for logging", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "ticketID", Value: id.String()})
	}
	oldValueStr := string(oldValueJSON)

	// Apply fields from the request DTO.
	// NOTE: Your DTO has *string for EventID, UserID, etc., but the entity has a specific type.
	// This part might need adjustment based on how you handle the conversion in the handler.
	// For now, we assume the handler will populate the status.
	if req.Status != nil {
		existingTicket.Status = entities.TicketStatus(*req.Status) // Assuming req.Status is a string that can be cast
	}
	if reason != "" {
		existingTicket.CancelReason = reason
		existingTicket.CancelledAt = new(time.Time)
		*existingTicket.CancelledAt = time.Now()
	}

	existingTicket.UpdatedAt = time.Now()
	appErr = s.ticketRepo.Update(ctx, existingTicket)
	if appErr != nil {
		return nil, appErr
	}

	// Log the activity with the old and new values.
	newValueJSON, err := json.Marshal(existingTicket)
	if err != nil {
		s.logger.Warn(ctx, "Failed to marshal new ticket value for logging", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "ticketID", Value: id.String()})
	}
	newValueStr := string(newValueJSON)
	ticketIDStr := existingTicket.ID.String()
	details := fmt.Sprintf("Ticket Status updated to '%s'. Ticket Code: %s", existingTicket.Status, existingTicket.TicketCode)
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, entities.ActionTypeUpdateTicketStatus, entities.ResourceTypeTicket, &ticketIDStr, &ipAddr, &details, &oldValueStr, &newValueStr)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log UpdateStatusTicket activity", logErr, logger.Field{Key: "ticketID", Value: existingTicket.ID.String()})
	}

	return existingTicket, nil
}

func (s *TicketServiceImpl) DeleteTicket(ctx context.Context, id utils.BinaryUUID, logCtx entities.ActivityLogContext) *exception.AppError {
	ticketToDelete, appErr := s.ticketRepo.GetByID(ctx, id)
	if appErr != nil {
		return appErr
	}
	oldValueJSON, _ := json.Marshal(ticketToDelete)
	oldValueStr := string(oldValueJSON)

	appErr = s.ticketRepo.Delete(ctx, id)
	if appErr != nil {
		return appErr
	}

	ticketIDStr := ticketToDelete.ID.String()
	// FIX: Corrected typo from ticketToDelete.ticketIDStr to ticketIDStr.
	details := fmt.Sprintf("Ticket '%s' deleted. Ticket Code: %s", ticketIDStr, ticketToDelete.TicketCode)
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, entities.ActionTypeDeleteTicket, entities.ResourceTypeTicket, &ticketIDStr, &ipAddr, &details, &oldValueStr, nil)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log DeleteTicket activity", logErr, logger.Field{Key: "ticketID", Value: ticketToDelete.ID.String()})
	}

	return nil
}

// FIX: ADDED this method to satisfy the contract.TicketService interface.
// This replaces the old `GetTicketByUser` which had an incorrect signature.
func (s *TicketServiceImpl) GetUserTickets(ctx context.Context, userID utils.BinaryUUID, filter contract.TicketFilter, logCtx entities.ActivityLogContext) ([]*entities.Ticket, int64, *exception.AppError) {
	// FIX: The repository's GetByUser method requires a filter. We pass it along.
	tickets, total, appErr := s.ticketRepo.GetByUser(ctx, userID, filter)
	if appErr != nil {
		s.logger.Error(ctx, "Failed to get user tickets from repository", appErr, logger.Field{Key: "userID", Value: userID.String()})
		return nil, 0, appErr
	}

	s.logger.Info(ctx, "Successfully retrieved user tickets",
		logger.Field{Key: "userID", Value: userID.String()},
		logger.Field{Key: "count", Value: len(tickets)})

	return tickets, total, nil
}

// Get All Tickets Service
func (s *TicketServiceImpl) GetAllTickets(ctx context.Context, filter contract.TicketFilter, logCtx entities.ActivityLogContext) ([]*entities.Ticket, int64, *exception.AppError) {
	// FIX: The repository method is named `GetAll`, not `GetAllTickets`.
	tickets, total, appErr := s.ticketRepo.GetAll(ctx, filter)
	if appErr != nil {
		s.logger.Error(ctx, "Failed to get all tickets from repository", appErr,
			logger.Field{Key: "filter", Value: filter})
		// FIX: Return the error correctly.
		return nil, 0, appErr
	}

	// FIX: REMOVED all DTO transformation logic. The service layer should return entities.
	// The handler is responsible for converting entities to DTOs for the web response.
	// This also fixes the `undefined: dto.TicketResponses` error.

	s.logger.Info(ctx, "Successfully retrieved all issued tickets",
		logger.Field{Key: "totalCount", Value: total},
		logger.Field{Key: "currentPage", Value: filter.Page})

	// FIX: Return the entities and total count directly from the repository call.
	return tickets, total, nil
}

// NOTE: The methods `GetTicketByUser` and `GetTicketByEvent` were removed as they are not part of the `contract.TicketService` interface.
// Their functionality is covered by `GetUserTickets` and `GetAllTickets` with appropriate filters.