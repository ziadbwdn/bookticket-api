package router

import (
	"root-app/internal/api/handler"
	"root-app/internal/logger"
	"root-app/internal/middleware"
	"root-app/internal/entities"

	"github.com/gin-gonic/gin"
)

func setupUserActivityRoutes(apiGroup *gin.RouterGroup, userActivityHandler *handler.UserActivityHandler, log logger.Logger) {
	activityGroup := apiGroup.Group("/activities")
	{
		activityGroup.POST("", userActivityHandler.LogActivity)
		activityGroup.GET("", middleware.AuthorizeRole(log, entities.RoleAdmin), userActivityHandler.ListActivities)
		activityGroup.GET("/summary/:userID", middleware.AuthorizeRole(log, entities.RoleAdmin), userActivityHandler.GetActivitySummary)
		activityGroup.GET("/alerts", middleware.AuthorizeRole(log, entities.RoleAdmin), userActivityHandler.GetSecurityAlerts)
		activityGroup.DELETE("", middleware.AuthorizeRole(log, entities.RoleAdmin), userActivityHandler.CleanOldActivities)
	}
}