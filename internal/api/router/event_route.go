package router

import (
	"root-app/internal/api/handler"
	"root-app/internal/entities"
	"root-app/internal/logger"
	"root-app/internal/middleware"

	"github.com/gin-gonic/gin"
)

// FIX: Added logger.Logger as a parameter.
func setupEventRoutes(r *gin.RouterGroup, eventHandler *handler.EventHandler, authMiddleware *middleware.AuthMiddleware, log logger.Logger) {
	events := r.Group("/events")
	{
		// FIX: Pass the 'log' variable to the middleware.
		events.POST("", middleware.AuthorizeRole(log, entities.RoleAdmin), eventHandler.CreateEvent)
		events.GET("/:id", eventHandler.GetEventByID)
		events.GET("", eventHandler.GetAllEvents)
		events.PUT("/:id", middleware.AuthorizeRole(log, entities.RoleAdmin), eventHandler.UpdateEvent)
		events.DELETE("/:id", middleware.AuthorizeRole(log, entities.RoleAdmin), eventHandler.DeleteEvent)
	}
}