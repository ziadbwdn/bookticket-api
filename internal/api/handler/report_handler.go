package handler

import (
	"net/http"
	"root-app/internal/contract"
	"root-app/internal/entities"
	"root-app/internal/exception"
	"root-app/internal/logger"
	"root-app/internal/utils"
	"root-app/pkg/gin_helper"
	"root-app/pkg/web_response"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService contract.ReportService
	logger        logger.Logger
}

func NewReportHandler(reportService contract.ReportService, logger logger.Logger) *ReportHandler {
	if reportService == nil {
		panic("reportService must not be nil for ReportHandler")
	}
	if logger == nil {
		panic("logger must not be nil for ReportHandler")
	}
	return &ReportHandler{
		reportService: reportService,
		logger:        logger,
	}
}

func (h *ReportHandler) GetSummaryReport(c *gin.Context) {
	userID, appErr := gin_helper.GetUserIDFromContext(c)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}
	userRole, appErr := gin_helper.GetUserRoleFromContext(c)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}
	if userRole != entities.RoleAdmin {
		appErr := exception.NewPermissionError("Only administrators can generate reports")
		web_response.HandleAppError(c, appErr)
		return
	}
	username, _ := gin_helper.GetUsernameFromContext(c)
	ipAddress := gin_helper.GetIPAddressFromContext(c)

	logCtx := entities.ActivityLogContext{
		// FIX: Convert the utils.BinaryUUID to a string here.
		UserID:    userID.String(),
		Username:  username,
		IPAddress: ipAddress,
	}

	summary, appErr := h.reportService.GenerateSummaryReport(c.Request.Context(), logCtx)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}
	h.logger.Info(c.Request.Context(), "Summary report generated successfully",
		logger.Field{Key: "userID", Value: userID.String()},
		logger.Field{Key: "totalTicketsSold", Value: summary.TotalTicketsSold})
	web_response.RespondWithSuccess(c, http.StatusOK, summary)
}

func (h *ReportHandler) GetTicketEventReport(c *gin.Context) {
	eventIDStr := c.Param("eventId")
	if eventIDStr == "" {
		appErr := exception.NewValidationError("Event ID is required")
		web_response.HandleAppError(c, appErr)
		return
	}
	// FIX: Corrected typo from StringToBinaryUUID to the existing ParseBinaryUUID function.
	eventID, err := utils.ParseBinaryUUID(eventIDStr)
	if err != nil {
		appErr := exception.NewValidationError("Invalid event ID format", err.Error())
		web_response.HandleAppError(c, appErr)
		return
	}
	userID, appErr := gin_helper.GetUserIDFromContext(c)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}
	userRole, appErr := gin_helper.GetUserRoleFromContext(c)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}
	if userRole != entities.RoleAdmin {
		appErr := exception.NewPermissionError("Only administrators can generate reports")
		web_response.HandleAppError(c, appErr)
		return
	}
	username, _ := gin_helper.GetUsernameFromContext(c)
	ipAddress := gin_helper.GetIPAddressFromContext(c)

	logCtx := entities.ActivityLogContext{
		// FIX: Convert the utils.BinaryUUID to a string here.
		UserID:    userID.String(),
		Username:  username,
		IPAddress: ipAddress,
	}
	report, appErr := h.reportService.GenerateTicketEventReport(c.Request.Context(), eventID, logCtx)
	if appErr != nil {
		web_response.HandleAppError(c, appErr)
		return
	}
	h.logger.Info(c.Request.Context(), "Ticket event report generated successfully",
		logger.Field{Key: "userID", Value: userID.String()},
		logger.Field{Key: "eventID", Value: eventID.String()},
		logger.Field{Key: "eventName", Value: report.EventName})
	web_response.RespondWithSuccess(c, http.StatusOK, report)
}