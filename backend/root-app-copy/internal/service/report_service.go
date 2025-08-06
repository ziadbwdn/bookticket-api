// internal/service/report_service.go
package service

import (
	"context"
	"fmt"
	"root-app/internal/contract"
	"root-app/internal/entities"
	"root-app/internal/exception"
	"root-app/internal/logger"
	"root-app/internal/utils"
)

// reportService implements the contract.ReportService interface.
type reportService struct {
	reportRepo      contract.ReportRepository
	activityService contract.UserActivityService
	logger          logger.Logger
}

// NewReportService creates and returns a new instance of ReportService.
func NewReportService(
	reportRepo contract.ReportRepository,
	activityService contract.UserActivityService,
	logger logger.Logger,
) contract.ReportService {
	if reportRepo == nil {
		panic("reportRepository must not be nil for ReportService")
	}
	if activityService == nil {
		panic("activityService must not be nil for ReportService")
	}
	if logger == nil {
		panic("logger must not be nil for ReportService")
	}
	return &reportService{
		reportRepo:      reportRepo,
		activityService: activityService,
		logger:          logger,
	}
}

// GenerateSummaryReport generates a comprehensive summary report and logs the activity synchronously.
func (s *reportService) GenerateSummaryReport(ctx context.Context, logCtx entities.ActivityLogContext) (*contract.ReportSummary, *exception.AppError) {
	op := "reportService.GenerateSummaryReport"
	s.logger.Info(ctx, "Attempting to generate summary report",
		logger.Field{Key: "userID", Value: logCtx.UserID},
		logger.Field{Key: "username", Value: logCtx.Username})

	summary, appErr := s.reportRepo.GetSummaryReport(ctx)
	if appErr != nil {
		s.logger.Error(ctx, "Failed to get summary report from repository", appErr.Err, logger.Field{Key: "operation", Value: op})
		// Synchronously log the activity of failed report generation
		details := fmt.Sprintf("Failed to generate summary report: %s", appErr.Message)
		if logErr := s.activityService.LogUserActivity(
			ctx,
			logCtx.UserID,
			logCtx.Username,
			entities.ActionTypeSummaryReport,
			entities.ResourceTypeReport,
			nil, // No specific resource ID for overall summary
			&logCtx.IPAddress,
			&details,
			nil,
			nil,
		); logErr != nil {
			s.logger.Error(ctx, "Failed to log user activity for failed summary report", logErr,
				logger.Field{Key: "userID", Value: logCtx.UserID},
				logger.Field{Key: "actionType", Value: entities.ActionTypeSummaryReport})
		}
		return nil, appErr
	}

	s.logger.Info(ctx, "Summary report generated successfully",
		logger.Field{Key: "totalTicketsSold", Value: summary.TotalTicketsSold},
		logger.Field{Key: "userID", Value: logCtx.UserID})

	// Synchronously log successful report generation activity
	details := "Successfully generated overall summary report."
	if logErr := s.activityService.LogUserActivity(
		ctx,
		logCtx.UserID,
		logCtx.Username,
		entities.ActionTypeSummaryReport,
		entities.ResourceTypeReport,
		nil, // No specific resource ID for overall summary
		&logCtx.IPAddress,
		&details,
		nil,
		nil,
	); logErr != nil {
		s.logger.Error(ctx, "Failed to log user activity for successful summary report", logErr,
			logger.Field{Key: "userID", Value: logCtx.UserID},
			logger.Field{Key: "actionType", Value: entities.ActionTypeSummaryReport})
	}

	return summary, nil
}

// GenerateTicketEventReport generates a detailed report for a specific event's tickets and logs the activity synchronously.
func (s *reportService) GenerateTicketEventReport(ctx context.Context, eventID utils.BinaryUUID, logCtx entities.ActivityLogContext) (*contract.TicketEventReport, *exception.AppError) {
	op := "reportService.GenerateTicketEventReport"
	s.logger.Info(ctx, "Attempting to generate ticket event report",
		logger.Field{Key: "eventID", Value: eventID.String()},
		logger.Field{Key: "userID", Value: logCtx.UserID})

	report, appErr := s.reportRepo.GetTicketEventReport(ctx, eventID)
	if appErr != nil {
		s.logger.Error(ctx, "Failed to get ticket event report from repository", appErr.Err,
			logger.Field{Key: "eventID", Value: eventID.String()}, logger.Field{Key: "operation", Value: op})
		// Synchronously log the activity of failed report generation
		eventIDStr := eventID.String()
		details := fmt.Sprintf("Failed to generate ticket event report for Event ID %s: %s", eventIDStr, appErr.Message)
		if logErr := s.activityService.LogUserActivity(
			ctx,
			logCtx.UserID,
			logCtx.Username,
			entities.ActionTypeTicketEventReport,
			entities.ResourceTypeReport,
			&eventIDStr, // Resource ID is the event ID
			&logCtx.IPAddress,
			&details,
			nil,
			nil,
		); logErr != nil {
			s.logger.Error(ctx, "Failed to log user activity for failed ticket event report", logErr,
				logger.Field{Key: "userID", Value: logCtx.UserID},
				logger.Field{Key: "actionType", Value: entities.ActionTypeTicketEventReport},
				logger.Field{Key: "eventID", Value: eventIDStr})
		}
		return nil, appErr
	}

	s.logger.Info(ctx, "Ticket event report generated successfully",
		logger.Field{Key: "eventID", Value: eventID.String()},
		logger.Field{Key: "ticketsSold", Value: report.TicketsSold},
		logger.Field{Key: "userID", Value: logCtx.UserID})

	// Synchronously log successful report generation activity
	eventIDStr := eventID.String()
	details := fmt.Sprintf("Successfully generated ticket event report for Event '%s' (ID: %s).", report.EventName, eventIDStr)
	if logErr := s.activityService.LogUserActivity(
		ctx,
		logCtx.UserID,
		logCtx.Username,
		entities.ActionTypeTicketEventReport,
		entities.ResourceTypeReport,
		&eventIDStr, // Resource ID is the event ID
		&logCtx.IPAddress,
		&details,
		nil,
		nil,
	); logErr != nil {
		s.logger.Error(ctx, "Failed to log user activity for successful ticket event report", logErr,
			logger.Field{Key: "userID", Value: logCtx.UserID},
			logger.Field{Key: "actionType", Value: entities.ActionTypeTicketEventReport},
			logger.Field{Key: "eventID", Value: eventIDStr})
	}

	return report, nil
}


