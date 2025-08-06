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

	"gorm.io/gorm"
)

// Event Service Struct implementation, containing of interfaces from repo
type EventServiceImpl struct {
	db              *gorm.DB // Add database connection
	eventRepo       contract.EventRepository
	activityService contract.UserActivityService
	logger          logger.Logger
}

// Event Service Constructor
func NewEventService(db *gorm.DB, eventRepo contract.EventRepository, activityService contract.UserActivityService, logger logger.Logger) contract.EventService {
	// panic checking
	if db == nil {
		panic("database connection must not be nil for EventService")
	}
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
		db:              db,
		eventRepo:       eventRepo,
		activityService: activityService,
		logger:          logger,
	}
}

func (s *EventServiceImpl) CreateEvent(ctx context.Context, event *entities.Event, logCtx entities.ActivityLogContext) (*entities.Event, *exception.AppError) {
	tx := s.db.Begin()
	if tx.Error != nil {
		s.logger.Error(ctx, "Failed to begin transaction for CreateEvent", tx.Error)
		return nil, exception.NewInternalError("failed to begin transaction", tx.Error)
	}

	event.ID = utils.NewBinaryUUID()
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}
	event.UpdatedAt = time.Now()

	// Pass the transaction context to the repository
	appErr := s.eventRepo.Create(ctx, event)
	if appErr != nil {
		tx.Rollback()
		return nil, appErr
	}

	// --- LOG USER ACTIVITY ---
	newValueJSON, _ := json.Marshal(event)
	newValueStr := string(newValueJSON)
	eventIDStr := event.ID.String()
	details := fmt.Sprintf("New Event '%s' created.", event.Name)
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, entities.ActionTypeCreateEvent, entities.ResourceTypeEvent, &eventIDStr, &ipAddr, &details, nil, &newValueStr)
	if logErr != nil {
		tx.Rollback() // Rollback if activity logging fails
		s.logger.Error(ctx, "Failed to log CreateEvent activity", logErr, logger.Field{Key: "eventID", Value: eventIDStr})
		return nil, exception.NewInternalError("failed to log activity", logErr)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // Ensure rollback if commit fails
		s.logger.Error(ctx, "Failed to commit transaction for CreateEvent", err)
		return nil, exception.NewInternalError("failed to commit transaction", err)
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
	tx := s.db.Begin()
	if tx.Error != nil {
		s.logger.Error(ctx, "Failed to begin transaction for UpdateEvent", tx.Error)
		return nil, exception.NewInternalError("failed to begin transaction", tx.Error)
	}

	existingEvent, appErr := s.eventRepo.GetByID(ctx, event.ID)
	if appErr != nil {
		tx.Rollback()
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
	if !event.StartDate.IsZero() {
		existingEvent.StartDate = event.StartDate
	}
	if !event.EndDate.IsZero() {
		existingEvent.EndDate = event.EndDate
	}
	if event.Capacity != 0 {
		existingEvent.Capacity = event.Capacity
	}

	existingEvent.UpdatedAt = time.Now()
	appErr = s.eventRepo.Update(ctx, existingEvent)
	if appErr != nil {
		tx.Rollback()
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
		tx.Rollback() // Rollback if activity logging fails
		s.logger.Error(ctx, "Failed to log UpdateEvent activity", logErr, logger.Field{Key: "eventID", Value: existingEvent.ID.String()})
		return nil, exception.NewInternalError("failed to log activity", logErr)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // Ensure rollback if commit fails
		s.logger.Error(ctx, "Failed to commit transaction for UpdateEvent", err)
		return nil, exception.NewInternalError("failed to commit transaction", err)
	}

	return existingEvent, nil
}

func (s *EventServiceImpl) DeleteEvent(ctx context.Context, id utils.BinaryUUID, logCtx entities.ActivityLogContext) *exception.AppError {
	tx := s.db.Begin()
	if tx.Error != nil {
		s.logger.Error(ctx, "Failed to begin transaction for DeleteEvent", tx.Error)
		return exception.NewInternalError("failed to begin transaction", tx.Error)
	}

	eventToDelete, appErr := s.eventRepo.GetByID(ctx, id)
	if appErr != nil {
		tx.Rollback()
		return appErr
	}
	oldValueJSON, _ := json.Marshal(eventToDelete)
	oldValueStr := string(oldValueJSON)

	appErr = s.eventRepo.Delete(ctx, id)
	if appErr != nil {
		tx.Rollback()
		return appErr
	}

	eventIDStr := eventToDelete.ID.String()
	details := fmt.Sprintf("Event '%s' deleted.", eventToDelete.Name)
	ipAddr := logCtx.IPAddress

	logErr := s.activityService.LogUserActivity(ctx, logCtx.UserID, logCtx.Username, entities.ActionTypeDeleteEvent, entities.ResourceTypeEvent, &eventIDStr, &ipAddr, &details, &oldValueStr, nil)
	if logErr != nil {
		tx.Rollback() // Rollback if activity logging fails
		s.logger.Error(ctx, "Failed to log DeleteEvent activity", logErr, logger.Field{Key: "eventID", Value: eventToDelete.ID.String()})
		return exception.NewInternalError("failed to log activity", logErr)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // Ensure rollback if commit fails
		s.logger.Error(ctx, "Failed to commit transaction for DeleteEvent", err)
		return exception.NewInternalError("failed to commit transaction", err)
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
