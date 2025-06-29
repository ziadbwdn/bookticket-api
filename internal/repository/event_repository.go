package repository

import (
	"context"
	"strings"
	"root-app/internal/contract"
	"root-app/internal/entities"
	"root-app/internal/exception"
	"root-app/internal/logger"
	"root-app/internal/utils"
	
	"gorm.io/gorm"
)

// compose event repo struct
type eventRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// repository constructor
func NewEventRepository(db *gorm.DB, logger logger.Logger) contract.EventRepository {
	return &eventRepository{
		db:     db,
		logger: logger,
	}
}

// Create Event repository logic
func (r *eventRepository) Create(ctx context.Context, event *entities.Event) *exception.AppError {
	if err := r.db.WithContext(ctx).Create(event).Error; err != nil {
		r.logger.Error(ctx, "Failed to create event", err, logger.Field{Key: "event_name", Value: event.Name})
		return exception.NewDatabaseError("create_event", err)
	}
	return nil
}

// Get Event by ID repository with filter
func (r *eventRepository) GetByID(ctx context.Context, id utils.BinaryUUID) (*entities.Event, *exception.AppError) {
	var event entities.Event
	if err := r.db.WithContext(ctx).Preload("Creator").First(&event, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("Event", id.String())
		}
		r.logger.Error(ctx, "Failed to get event by ID", err, logger.Field{Key: "event_id", Value: id.String()})
		return nil, exception.NewDatabaseError("get_event_by_id", err)
	}
	return &event, nil
}

// Get All Repository logic with filter
func (r *eventRepository) GetAll(ctx context.Context, filter contract.EventFilter) ([]*entities.Event, int64, *exception.AppError) {
	var events []*entities.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Event{})

	// Apply filters
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		searchTerm := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error(ctx, "Failed to count events", err)
		return nil, 0, exception.NewDatabaseError("count_events", err)
	}

	// Apply pagination
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	// Execute query
	if err := query.Preload("Creator").Order("created_at DESC").Find(&events).Error; err != nil {
		r.logger.Error(ctx, "Failed to get events", err)
		return nil, 0, exception.NewDatabaseError("get_events", err)
	}

	return events, total, nil
}

func (r *eventRepository) Update(ctx context.Context, event *entities.Event) *exception.AppError {
	if err := r.db.WithContext(ctx).Save(event).Error; err != nil {
		r.logger.Error(ctx, "Failed to update event", err, logger.Field{Key: "event_id", Value: event.ID.String()})
		return exception.NewDatabaseError("update_event", err)
	}
	return nil
}

func (r *eventRepository) Delete(ctx context.Context, id utils.BinaryUUID) *exception.AppError {
	if err := r.db.WithContext(ctx).Delete(&entities.Event{}, "id = ?", id).Error; err != nil {
		r.logger.Error(ctx, "Failed to delete event", err, logger.Field{Key: "event_id", Value: id.String()})
		return exception.NewDatabaseError("delete_event", err)
	}
	return nil
}

func (r *eventRepository) CheckNameExists(ctx context.Context, name string, excludeID *utils.BinaryUUID) (bool, *exception.AppError) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Event{}).Where("name = ?", name)
	
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	
	if err := query.Count(&count).Error; err != nil {
		r.logger.Error(ctx, "Failed to check event name existence", err, logger.Field{Key: "event_name", Value: name})
		return false, exception.NewDatabaseError("check_event_name", err)
	}
	
	return count > 0, nil
}

func (r *eventRepository) GetAvailableCapacity(ctx context.Context, eventID utils.BinaryUUID) (int, *exception.AppError) {
	var event entities.Event
	if err := r.db.WithContext(ctx).First(&event, "id = ?", eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, exception.NewNotFoundError("Event", eventID.String())
		}
		return 0, exception.NewDatabaseError("get_event_capacity", err)
	}

	var soldTickets int64
	if err := r.db.WithContext(ctx).Model(&entities.Ticket{}).
		Where("event_id = ? AND status IN (?)", eventID, []string{"active", "used"}).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&soldTickets).Error; err != nil {
		return 0, exception.NewDatabaseError("count_sold_tickets", err)
	}

	availableCapacity := event.Capacity - int(soldTickets)
	if availableCapacity < 0 {
		availableCapacity = 0
	}

	return availableCapacity, nil
}