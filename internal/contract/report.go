package contract

import (
	"context"
	"root-app/internal/exception"
	"root-app/internal/entities"
	"root-app/internal/utils"
	"time"
)

type ReportSummary struct {
	TotalTicketsSold int                `json:"totalTicketsSold"`
	TotalRevenue     *utils.GormDecimal `json:"totalRevenue"`
	TotalEvents      int                `json:"totalEvents"`
	ActiveEvents     int                `json:"activeEvents"`
	GeneratedAt      time.Time          `json:"generatedAt"`
}

type TicketEventReport struct {
	EventID          utils.BinaryUUID   `json:"eventId"`
	EventName        string             `json:"eventName"`
	EventCategory    string             `json:"eventCategory"`
	EventCapacity    int                `json:"eventCapacity"`
	TicketsSold      int                `json:"ticketsSold"`
	AvailableTickets int                `json:"availableTickets"`
	Revenue          *utils.GormDecimal `json:"revenue"`
	Status           string             `json:"status"`
	GeneratedAt      time.Time          `json:"generatedAt"`
}

type ReportRepository interface {
	GetSummaryReport(ctx context.Context) (*ReportSummary, *exception.AppError)
	GetTicketEventReport(ctx context.Context, eventID utils.BinaryUUID) (*TicketEventReport, *exception.AppError)
	GetRevenueByDateRange(ctx context.Context, startDate, endDate time.Time) (*utils.GormDecimal, *exception.AppError) // optional 
}

type ReportService interface {
	GenerateSummaryReport(ctx context.Context, logCtx entities.ActivityLogContext) (*ReportSummary, *exception.AppError)
	GenerateTicketEventReport(ctx context.Context, eventID utils.BinaryUUID, logCtx entities.ActivityLogContext) (*TicketEventReport, *exception.AppError)
}