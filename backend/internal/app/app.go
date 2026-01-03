package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/cache"
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/router"

	"gorm.io/gorm"
)

type App struct {
	config *config.Config
	db     *gorm.DB
	cache  *cache.RedisCache
	server *http.Server
}

func New() (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	app := &App{
		config: cfg,
	}

	// Initialize database
	if err := app.initDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize cache (optional)
	if err := app.initCache(); err != nil {
		log.Printf("Warning: Failed to initialize cache: %v", err)
	}

	// Initialize server
	app.initServer()

	return app, nil
}

func (a *App) initDatabase() error {
	dbConfig := database.Config{
		Host:     a.config.Database.Host,
		Port:     a.config.Database.Port,
		User:     a.config.Database.User,
		Password: a.config.Database.Password,
		DBName:   a.config.Database.DBName,
		SSLMode:  a.config.Database.SSLMode,
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		return err
	}

	// Run auto migrations
	if err := database.AutoMigrate(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	a.db = db
	log.Println("✓ Database initialized successfully")
	return nil
}

func (a *App) initCache() error {
	if !a.config.Redis.Enabled {
		log.Println("ℹ Redis cache is disabled")
		return nil
	}

	redisCache, err := cache.NewRedisCache(
		a.config.Redis.Addr,
		a.config.Redis.Password,
		a.config.Redis.DB,
	)
	if err != nil {
		return err
	}

	a.cache = redisCache
	log.Println("✓ Redis cache initialized successfully")
	return nil
}

func (a *App) initServer() {
	container := NewContainer(a.db, a.cache)

	r := router.New(container.Handlers())
	ginRouter := r.Setup(a.config.IsDevelopment())

	a.server = &http.Server{
		Addr:    ":" + a.config.Server.Port,
		Handler: ginRouter,
	}
}

func (a *App) Run() error {
	// Channel to listen for errors
	serverErrors := make(chan error, 1)

	// Start HTTP server
	go func() {
		log.Printf("🚀 Server starting on port %s", a.config.Server.Port)
		log.Printf("📝 Environment: %s", a.config.Server.Env)
		log.Printf("🔗 Health check: http://localhost:%s/health", a.config.Server.Port)
		log.Printf("🔗 API docs: http://localhost:%s/", a.config.Server.Port)

		serverErrors <- a.server.ListenAndServe()
	}()

	// Channel to listen for interrupt signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or server error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Printf("🛑 Received signal: %v. Starting graceful shutdown...", sig)

		// Create context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Attempt graceful shutdown
		if err := a.Shutdown(ctx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}

		log.Println("✓ Server stopped gracefully")
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	// Shutdown HTTP server
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	// Close database connection
	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Warning: Failed to close database: %v", err)
			} else {
				log.Println("✓ Database connection closed")
			}
		}
	}

	// Close Redis connection
	if a.cache != nil {
		if err := a.cache.Close(); err != nil {
			log.Printf("Warning: Failed to close Redis: %v", err)
		} else {
			log.Println("✓ Redis connection closed")
		}
	}

	return nil
}

func (a *App) GetDB() *gorm.DB {
	return a.db
}

func (a *App) GetCache() *cache.RedisCache {
	return a.cache
}

func (a *App) GetConfig() *config.Config {
	return a.config
}
