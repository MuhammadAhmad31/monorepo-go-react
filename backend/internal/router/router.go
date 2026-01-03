package router

import (
	"backend/internal/generated"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	handler *handlers.CombinedHandler
}

func New(handler *handlers.CombinedHandler) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) Setup(isDevelopment bool) *gin.Engine {
	if !isDevelopment {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())
	router.Use(gin.Logger())
	router.Use(middleware.RateLimit(100, time.Minute))

	// Health check and welcome (public endpoints)
	router.GET("/health", r.healthCheck)
	router.GET("/", r.welcome)

	// API v1 group with RBAC middleware
	v1 := router.Group("/api/v1")

	// Apply OpenAPI-based RBAC middleware
	v1.Use(middleware.OpenAPISecurityMiddleware())

	// Register oapi-codegen generated handlers
	// Security is now handled by OpenAPISecurityMiddleware
	generated.RegisterHandlers(v1, r.handler)

	return router
}

func (r *Router) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": "backend-api",
	})
}

func (r *Router) welcome(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to Backend API",
		"version": "1.0.0",
		"endpoints": gin.H{
			"health":   "/health",
			"auth":     "/api/v1/auth",
			"users":    "/api/v1/users",
			"products": "/api/v1/products",
		},
	})
}
