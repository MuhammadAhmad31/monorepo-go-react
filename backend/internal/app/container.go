package app

import (
	"backend/internal/cache"
	"backend/internal/handlers"
	"backend/internal/repository"
	"backend/internal/service"

	"gorm.io/gorm"
)

type Container struct {
	UserHandler    *handlers.UserHandler
	ProductHandler *handlers.ProductHandler
	AuthHandler    *handlers.AuthHandler
}

func NewContainer(db *gorm.DB, cache *cache.RedisCache) *Container {
	// repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)

	// services
	userService := service.NewUserService(userRepo, cache)
	productService := service.NewProductService(productRepo, cache)
	authService := service.NewAuthService(userRepo)

	// handlers
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	authHandler := handlers.NewAuthHandler(authService)

	return &Container{
		UserHandler:    userHandler,
		ProductHandler: productHandler,
		AuthHandler:    authHandler,
	}
}

func (c *Container) Handlers() *handlers.CombinedHandler {
	return &handlers.CombinedHandler{
		UserHandler:    c.UserHandler,
		ProductHandler: c.ProductHandler,
		AuthHandler:    c.AuthHandler,
	}
}
