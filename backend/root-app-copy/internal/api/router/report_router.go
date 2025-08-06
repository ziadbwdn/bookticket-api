package router

import (
	"root-app/internal/api/handler"
	"root-app/internal/entities"
	"root-app/internal/logger"
	"root-app/internal/middleware"

	"github.com/gin-gonic/gin"
)

// FIX: Added logger.Logger as a parameter.
func setupReportRoutes(r *gin.RouterGroup, reportHandler *handler.ReportHandler, authMiddleware *middleware.AuthMiddleware, log logger.Logger) {
	reports := r.Group("/reports")
	{
		// FIX: Pass the 'log' variable to the middleware.
		reports.GET("/summary", middleware.AuthorizeRole(log, entities.RoleAdmin), reportHandler.GetSummaryReport)
		reports.GET("/events/:id", middleware.AuthorizeRole(log, entities.RoleAdmin), reportHandler.GetTicketEventReport)
	}
}