package router

import (
	"root-app/internal/config"
	"root-app/internal/api/handler"
	"root-app/internal/logger"
	"root-app/internal/middleware"
	"root-app/internal/repository"
	"root-app/internal/service"
	"root-app/pkg/web_response"
	"net/http"
	"time"
	
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// RouterConfig holds the core dependencies needed to build the API layer.
type RouterConfig struct {
	DB     *gorm.DB
	Cfg    *config.Config
	Logger logger.Logger
}

func SetupRouter(cfg *RouterConfig) *gin.Engine {
	// --- 1. Initialize Dependencies ---
	userRepo := repository.NewUserRepository(cfg.DB)
	eventRepo := repository.NewEventRepository(cfg.DB, cfg.Logger)
	ticketRepo := repository.NewTicketRepository(cfg.DB, cfg.Logger)
	reportRepo := repository.NewReportRepository(cfg.DB, cfg.Logger)
	userActivityRepo := repository.NewGormUserActivityRepository(cfg.DB, cfg.Logger)

	userActivityService := service.NewUserActivityService(cfg.DB, userActivityRepo, cfg.Logger)
	// Assuming you fixed NewAuthService to return contract.UserService
	authService := service.NewAuthService(userRepo, userActivityService, cfg.Cfg.JWTSecret)
	eventService := service.NewEventService(cfg.DB, eventRepo, userActivityService, cfg.Logger)
	ticketService := service.NewTicketService(cfg.DB, ticketRepo, eventRepo, userActivityService, cfg.Logger)
	reportService := service.NewReportService(reportRepo, userActivityService, cfg.Logger)

	authHandler := handler.NewAuthHandler(authService)
	eventHandler := handler.NewEventHandler(eventService, authService)
	ticketHandler := handler.NewTicketHandler(ticketService, authService)
	reportHandler := handler.NewReportHandler(reportService, cfg.Logger)
	userActivityHandler := handler.NewUserActivityHandler(userActivityService, cfg.Logger)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	// --- 2. Setup Router Engine ---
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", func(c *gin.Context) {
		web_response.RespondWithSuccess(c, http.StatusOK, gin.H{"status": "ok"})
	})

	// --- 3. Setup API Route Groups ---
	api := router.Group("/api")
	{
		// Public routes for authentication (login, register, etc.)
		// FIX: Removed the logger argument, as setupAuthRoutes does not require it.
		setupAuthRoutes(api.Group("/auth"), authHandler, authMiddleware)

		// public route for event, remove if not needed
		/*
		api.GET("/events", eventHandler.GetAllEvents)
		api.GET("/events/:id", eventHandler.GetEventByID)
		*/

		// Protected routes are grouped under the main authentication middleware
		protectedAPI := api.Group("/")
		protectedAPI.Use(authMiddleware.Handle())
		{
			// Call the setup functions for each domain, passing dependencies as needed.
			setupEventRoutes(protectedAPI, eventHandler, authMiddleware, cfg.Logger)
			setupTicketRoutes(protectedAPI, ticketHandler, authMiddleware, cfg.Logger)
			setupReportRoutes(protectedAPI, reportHandler, authMiddleware, cfg.Logger)
			
			// FIX: Added the call to setupUserActivityRoutes to use the handler.
			// Note that it does not take the authMiddleware parameter, as per your file.
			setupUserActivityRoutes(protectedAPI, userActivityHandler, cfg.Logger)
		}
	}

	return router
}