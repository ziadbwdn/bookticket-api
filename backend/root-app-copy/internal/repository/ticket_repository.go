package repository

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"root-app/internal/contract"
	"root-app/internal/entities"
	"root-app/internal/exception"
	"root-app/internal/logger"
	"root-app/internal/utils"
	
	"gorm.io/gorm"
)

type ticketRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewTicketRepository(db *gorm.DB, logger logger.Logger) contract.TicketRepository {
	return &ticketRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *entities.Ticket) *exception.AppError {
	if err := r.db.WithContext(ctx).Create(ticket).Error; err != nil {
		r.logger.Error(ctx, "Failed to create ticket", err, logger.Field{Key: "ticket_code", Value: ticket.TicketCode})
		return exception.NewDatabaseError("create_ticket", err)
	}
	return nil
}

func (r *ticketRepository) GetByID(ctx context.Context, id utils.BinaryUUID) (*entities.Ticket, *exception.AppError) {
	var ticket entities.Ticket
	if err := r.db.WithContext(ctx).Preload("Event").Preload("User").First(&ticket, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("Ticket", id.String())
		}
		r.logger.Error(ctx, "Failed to get ticket by ID", err, logger.Field{Key: "ticket_id", Value: id.String()})
		return nil, exception.NewDatabaseError("get_ticket_by_id", err)
	}
	return &ticket, nil
}

func (r *ticketRepository) GetAll(ctx context.Context, filter contract.TicketFilter) ([]*entities.Ticket, int64, *exception.AppError) {
	var tickets []*entities.Ticket
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Ticket{})

	// Apply filters
	if filter.EventID != (utils.BinaryUUID{}) {
		query = query.Where("event_id = ?", filter.EventID)
	}
	if filter.UserID != (utils.BinaryUUID{}) {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error(ctx, "Failed to count tickets", err)
		return nil, 0, exception.NewDatabaseError("count_tickets", err)
	}

	// Apply pagination
	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	// Execute query
	if err := query.Preload("Event").Preload("User").Order("created_at DESC").Find(&tickets).Error; err != nil {
		r.logger.Error(ctx, "Failed to get tickets", err)
		return nil, 0, exception.NewDatabaseError("get_tickets", err)
	}

	return tickets, total, nil
}

func (r *ticketRepository) Update(ctx context.Context, ticket *entities.Ticket) *exception.AppError {
	if err := r.db.WithContext(ctx).Save(ticket).Error; err != nil {
		r.logger.Error(ctx, "Failed to update ticket", err, logger.Field{Key: "ticket_id", Value: ticket.ID.String()})
		return exception.NewDatabaseError("update_ticket", err)
	}
	return nil
}

// Delete Ticket
func (r *ticketRepository) Delete(ctx context.Context, id utils.BinaryUUID) *exception.AppError {
	if err := r.db.WithContext(ctx).Delete(&entities.Ticket{}, "id = ?", id).Error; err != nil {
		r.logger.Error(ctx, "Failed to delete ticket", err, logger.Field{Key: "ticket_id", Value: id.String()})
		return exception.NewDatabaseError("delete_ticket", err)
	}
	return nil
}

// Get by user ID
func (r *ticketRepository) GetByUser(ctx context.Context, userID utils.BinaryUUID, filter contract.TicketFilter) ([]*entities.Ticket, int64, *exception.AppError) {
	filter.UserID = userID
	return r.GetAll(ctx, filter)
}

// Count ticket by Event
func (r *ticketRepository) CountTicketsByEvent(ctx context.Context, eventID utils.BinaryUUID) (int, *exception.AppError) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&entities.Ticket{}).
		Where("event_id = ? AND status IN (?)", eventID, []string{"active", "used"}).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&total).Error; err != nil {
		r.logger.Error(ctx, "Failed to count tickets by event", err, logger.Field{Key: "event_id", Value: eventID.String()})
		return 0, exception.NewDatabaseError("count_tickets_by_event", err)
	}
	return int(total), nil
}

// Get ticket by event
func (r *ticketRepository) GetTicketsByEvent(ctx context.Context, eventID utils.BinaryUUID) ([]*entities.Ticket, *exception.AppError) {
	var tickets []*entities.Ticket
	if err := r.db.WithContext(ctx).Preload("User").
		Where("event_id = ?", eventID).
		Order("created_at DESC").
		Find(&tickets).Error; err != nil {
		r.logger.Error(ctx, "Failed to get tickets by event", err, logger.Field{Key: "event_id", Value: eventID.String()})
		return nil, exception.NewDatabaseError("get_tickets_by_event", err)
	}
	return tickets, nil
}

// generate Unique Code
func (r *ticketRepository) GenerateUniqueTicketCode(ctx context.Context) (string, *exception.AppError) {
	const maxAttempts = 10
	rand.Seed(time.Now().UnixNano())

	for attempt := 0; attempt < maxAttempts; attempt++ {
		// Generate ticket code: TKT-YYYYMMDD-XXXXXX
		now := time.Now()
		dateStr := now.Format("20060102")
		randomStr := fmt.Sprintf("%06d", rand.Intn(1000000))
		ticketCode := fmt.Sprintf("TKT-%s-%s", dateStr, randomStr)

		// Check if code exists
		var count int64
		if err := r.db.WithContext(ctx).Model(&entities.Ticket{}).
			Where("ticket_code = ?", ticketCode).
			Count(&count).Error; err != nil {
			r.logger.Error(ctx, "Failed to check ticket code uniqueness", err, logger.Field{Key: "ticket_code", Value: ticketCode})
			return "", exception.NewDatabaseError("check_ticket_code", err)
		}

		if count == 0 {
			return ticketCode, nil
		}
	}

	return "", exception.NewInternalError("generate_ticket_code", fmt.Errorf("failed to generate unique ticket code after %d attempts", maxAttempts))
}