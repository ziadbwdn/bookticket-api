package router

import (
	"root-app/internal/api/handler"
	"root-app/internal/entities"
	"root-app/internal/logger"
	"root-app/internal/middleware"

	"github.com/gin-gonic/gin"
)

// FIX: Added logger.Logger as a parameter to the function signature.
func setupTicketRoutes(r *gin.RouterGroup, ticketHandler *handler.TicketHandler, authMiddleware *middleware.AuthMiddleware, log logger.Logger) {
	tickets := r.Group("/tickets")
	{
		// All routes in this group will require authentication because they are set up
		// under the protected group in router.go.

		tickets.POST("", ticketHandler.PurchaseTicket)
		tickets.GET("/:id", ticketHandler.GetTicketByID)
		tickets.PATCH("/:id/status", ticketHandler.UpdateTicketStatus)

		// These routes require an additional role check for Admin.
		// FIX: Corrected undefined 'RoleAdmin' to 'entities.RoleAdmin' and passed the 'log' variable.
		tickets.GET("", middleware.AuthorizeRole(log, entities.RoleAdmin), ticketHandler.GetAllTickets)
		// FIX: Corrected method call and used the consistent AuthorizeRole pattern.
		tickets.DELETE("/:id", middleware.AuthorizeRole(log, entities.RoleAdmin), ticketHandler.DeleteTicket)
	}
}