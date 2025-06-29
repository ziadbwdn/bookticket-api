package service

import (
	"context"
	"encoding/json"
	"fmt"
	"root-app/internal/contract"
	"root-app/internal/entities"
	"root-app/internal/exception"
	"root-app/internal/logger"
	"root-app/internal/utils"
	"time"
)

// Event Service Struct implementation, containing of interfaces from repo
type EventServiceImpl struct {
	eventRepo       contract.EventRepository
	activityService contract.UserActivityService
	logger          logger.Logger
}

// Event Service Constructor
func NewEventService(eventRepo contract.EventRepository, activityService contract.UserActivityService, logger logger.Logger) contract.EventService {
	// panic checking
	if eventRepo == nil {
		panic("eventRepo must not be nil for EventService")
	}
	if activityService == nil {
		panic("activityService must not be nil for EventService")
	}
	if logger == nil {
		panic("logger must not be nil for EventService")
	}
	return &EventServiceImpl{
		eventRepo:       eventRepo,
		activityService: activityService,
		logger:          logger,
	}
}

func (s *EventServiceImpl) CreateEvent(ctx context.Context, event *entities.Event, logCtx entities.ActivityLogContext) (*entities.Event, *exception.AppError) {
	event.ID = utils.NewBinaryUUID()
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}
	event.UpdatedAt = time.Now()

	appErr := s.eventRepo.Create(ctx, event)
	if appErr != nil {
		return nil, appErr
	}

	// --- LOG USER ACTIVITY ---
	newValueJSON, _ := json.Marshal(event)
	newValueStr := string(newValueJSON)
	eventIDStr := event.ID.String()
	details := fmt.Sprintf("New Event '%s' created.", event.Name)
	// FIX: The ActivityLogContext UserID is a string, but the LogUserActivity service expects a BinaryUUID.
	// This needs to be consistent. Assuming LogUserActivity takes a string for now as per handler code.
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, entities.ActionTypeCreateEvent, entities.ResourceTypeEvent, &eventIDStr, &ipAddr, &details, nil, &newValueStr)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log CreateEvent activity", logErr, logger.Field{Key: "eventID", Value: eventIDStr})
	}

	return event, nil
}

func (s *EventServiceImpl) GetEventByID(
	ctx context.Context,
	id utils.BinaryUUID,
) (*entities.Event, *exception.AppError) {
	event, appErr := s.eventRepo.GetByID(ctx, id)
	if appErr != nil {
		return nil, appErr
	}

	return event, nil
}

func (s *EventServiceImpl) UpdateEvent(ctx context.Context, event *entities.Event, logCtx entities.ActivityLogContext) (*entities.Event, *exception.AppError) {
	existingEvent, appErr := s.eventRepo.GetByID(ctx, event.ID)
	if appErr != nil {
		// FIX: Corrected return signature. Must return (*entities.Event, *exception.AppError).
		return nil, appErr
	}

	oldValueJSON, err := json.Marshal(existingEvent)
	if err != nil {
		s.logger.Warn(ctx, "Failed to marshal old event value for logging", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "eventID", Value: event.ID.String()})
	}
	oldValueStr := string(oldValueJSON)

	if event.Name != "" {
		existingEvent.Name = event.Name
	}
	if event.Description != "" {
		existingEvent.Description = event.Description
	}
	if event.Category != "" {
		existingEvent.Category = event.Category
	}
	if event.Venue != "" {
		existingEvent.Venue = event.Venue
	}
	// FIX: The check should be `!event.StartDate.IsZero()` to update only if a new date is provided.
	if !event.StartDate.IsZero() {
		existingEvent.StartDate = event.StartDate
	}
	if !event.EndDate.IsZero() {
		existingEvent.EndDate = event.EndDate
	}
	if event.Capacity != 0 {
		// FIX: This was a typo, it was updating Venue instead of Capacity.
		existingEvent.Capacity = event.Capacity
	}

	existingEvent.UpdatedAt = time.Now()
	// FIX: The repository's Update method only returns an error. Do not reassign existingEvent here.
	appErr = s.eventRepo.Update(ctx, existingEvent)
	if appErr != nil {
		// FIX: Corrected return signature.
		return nil, appErr
	}

	// Log the activity
	newValueJSON, err := json.Marshal(existingEvent)
	if err != nil {
		s.logger.Warn(ctx, "Failed to marshal new event value for logging", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "eventID", Value: existingEvent.ID.String()})
	}
	newValueStr := string(newValueJSON)
	eventIDStr := existingEvent.ID.String()
	details := fmt.Sprintf("Event '%s' updated.", existingEvent.Name)
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, entities.ActionTypeUpdateEvent, entities.ResourceTypeEvent, &eventIDStr, &ipAddr, &details, &oldValueStr, &newValueStr)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log UpdateEvent activity", logErr, logger.Field{Key: "eventID", Value: existingEvent.ID.String()})
	}

	return existingEvent, nil
}

func (s *EventServiceImpl) DeleteEvent(ctx context.Context, id utils.BinaryUUID, logCtx entities.ActivityLogContext) *exception.AppError {
	eventToDelete, appErr := s.eventRepo.GetByID(ctx, id)
	if appErr != nil {
		return appErr
	}
	oldValueJSON, _ := json.Marshal(eventToDelete)
	oldValueStr := string(oldValueJSON)

	appErr = s.eventRepo.Delete(ctx, id)
	if appErr != nil {
		return appErr
	}

	eventIDStr := eventToDelete.ID.String()
	details := fmt.Sprintf("Event '%s' deleted.", eventToDelete.Name)
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, entities.ActionTypeDeleteEvent, entities.ResourceTypeEvent, &eventIDStr, &ipAddr, &details, &oldValueStr, nil)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log DeleteEvent activity", logErr, logger.Field{Key: "eventID", Value: eventToDelete.ID.String()})
	}

	return nil
}

func (s *EventServiceImpl) GetAllEvents(ctx context.Context,
	filter contract.EventFilter,
	logCtx entities.ActivityLogContext) ([]*entities.Event, int64, *exception.AppError) {
	events, total, appErr := s.eventRepo.GetAll(ctx, filter)
	if appErr != nil {
		return nil, 0, appErr
	}

	details := "Viewed all events."
	ipAddr := logCtx.IPAddress
	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, entities.ActionTypeGetAllEvents, entities.ResourceTypeEvent, nil, &ipAddr, &details, nil, nil)
	if logErr != nil {
		s.logger.Error(ctx, "Failed to log GetAllEvents activity", logErr)
	}

	return events, total, nil
}