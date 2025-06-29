// internal/repository/report_repository.go
package repository

import (
	"context"
	"root-app/internal/contract"
	"root-app/internal/entities"
	"root-app/internal/exception"
	"root-app/internal/logger"
	"root-app/internal/utils"
	"time"
	
	"gorm.io/gorm"
)

type reportRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewReportRepository(db *gorm.DB, logger logger.Logger) contract.ReportRepository {
	return &reportRepository{
		db:     db,
		logger: logger,
	}
}

func (r *reportRepository) GetSummaryReport(ctx context.Context) (*contract.ReportSummary, *exception.AppError) {
	var summary contract.ReportSummary
	
	// Get total tickets sold
	var totalTicketsSold int64
	if err := r.db.WithContext(ctx).Model(&entities.Ticket{}).
		Where("status IN (?)", []string{"active", "used"}).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&totalTicketsSold).Error; err != nil {
		r.logger.Error(ctx, "Failed to get total tickets sold", err)
		return nil, exception.NewDatabaseError("get_total_tickets_sold", err)
	}
	summary.TotalTicketsSold = int(totalTicketsSold)

	// Get total revenue
	var totalRevenue string
	if err := r.db.WithContext(ctx).Model(&entities.Ticket{}).
		Where("status IN (?)", []string{"active", "used"}).
		Select("COALESCE(SUM(CAST(total_price AS DECIMAL(10,2))), 0)").
		Scan(&totalRevenue).Error; err != nil {
		r.logger.Error(ctx, "Failed to get total revenue", err)
		return nil, exception.NewDatabaseError("get_total_revenue", err)
	}
	
	revenueDecimal, appErr := utils.StringToGormDecimal(totalRevenue)
	if appErr != nil {
		return nil, appErr
	}
	summary.TotalRevenue = revenueDecimal

	// Get total events
	var totalEvents int64
	if err := r.db.WithContext(ctx).Model(&entities.Event{}).Count(&totalEvents).Error; err != nil {
		r.logger.Error(ctx, "Failed to get total events", err)
		return nil, exception.NewDatabaseError("get_total_events", err)
	}
	summary.TotalEvents = int(totalEvents)

	// Get active events
	var activeEvents int64
	if err := r.db.WithContext(ctx).Model(&entities.Event{}).
		Where("status = ?", entities.EventStatusActive).
		Count(&activeEvents).Error; err != nil {
		r.logger.Error(ctx, "Failed to get active events", err)
		return nil, exception.NewDatabaseError("get_active_events", err)
	}
	summary.ActiveEvents = int(activeEvents)

	summary.GeneratedAt = time.Now()
	return &summary, nil
}

func (r *reportRepository) GetTicketEventReport(ctx context.Context, eventID utils.BinaryUUID) (*contract.TicketEventReport, *exception.AppError) {
	var event entities.Event
	if err := r.db.WithContext(ctx).First(&event, "id = ?", eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundError("Event", eventID.String())
		}
		r.logger.Error(ctx, "Failed to get event for report", err, logger.Field{Key: "event_id", Value: eventID.String()})
		return nil, exception.NewDatabaseError("get_event_for_report", err)
	}

	var report contract.TicketEventReport
	report.EventID = event.ID
	report.EventName = event.Name
	report.EventCategory = event.Category
	report.EventCapacity = event.Capacity
	report.Status = string(event.Status)

	// Get tickets sold
	var ticketsSold int64
	if err := r.db.WithContext(ctx).Model(&entities.Ticket{}).
		Where("event_id = ? AND status IN (?)", eventID, []string{"active", "used"}).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&ticketsSold).Error; err != nil {
		r.logger.Error(ctx, "Failed to get tickets sold for event", err, logger.Field{Key: "event_id", Value: eventID.String()})
		return nil, exception.NewDatabaseError("get_tickets_sold_for_event", err)
	}
	report.TicketsSold = int(ticketsSold)
	report.AvailableTickets = event.Capacity - int(ticketsSold)

	// Get revenue for event
	var revenue string
	if err := r.db.WithContext(ctx).Model(&entities.Ticket{}).
		Where("event_id = ? AND status IN (?)", eventID, []string{"active", "used"}).
		Select("COALESCE(SUM(CAST(total_price AS DECIMAL(10,2))), 0)").
		Scan(&revenue).Error; err != nil {
		r.logger.Error(ctx, "Failed to get revenue for event", err, logger.Field{Key: "event_id", Value: eventID.String()})
		return nil, exception.NewDatabaseError("get_revenue_for_event", err)
	}

	revenueDecimal, appErr := utils.StringToGormDecimal(revenue)
	if appErr != nil {
		return nil, appErr
	}
	report.Revenue = revenueDecimal

	report.GeneratedAt = time.Now()
	return &report, nil
}

func (r *reportRepository) GetRevenueByDateRange(ctx context.Context, startDate, endDate time.Time) (*utils.GormDecimal, *exception.AppError) {
	var revenue string
	if err := r.db.WithContext(ctx).Model(&entities.Ticket{}).
		Where("status IN (?) AND created_at BETWEEN ? AND ?", []string{"active", "used"}, startDate, endDate).
		Select("COALESCE(SUM(CAST(total_price AS DECIMAL(10,2))), 0)").
		Scan(&revenue).Error; err != nil {
		r.logger.Error(ctx, "Failed to get revenue by date range", err)
		return nil, exception.NewDatabaseError("get_revenue_by_date_range", err)
	}

	return utils.StringToGormDecimal(revenue)
}