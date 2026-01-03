package service

import (
	"backend/internal/cache"
	"backend/internal/generated"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"fmt"
	"time"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProduct(ctx context.Context, id generated.IdParam) (*models.Product, error)
	ListProducts(ctx context.Context, page, perPage int) ([]models.Product, int64, error)
	UpdateProduct(ctx context.Context, id generated.IdParam, product *models.Product) error
	DeleteProduct(ctx context.Context, id generated.IdParam) error
}

type productService struct {
	repo  repository.ProductRepository
	cache *cache.RedisCache
}

func NewProductService(repo repository.ProductRepository, cache *cache.RedisCache) ProductService {
	return &productService{
		repo:  repo,
		cache: cache,
	}
}

func (s *productService) CreateProduct(ctx context.Context, product *models.Product) error {
	if err := s.repo.Create(ctx, product); err != nil {
		return err
	}

	// Invalidate products list cache
	if s.cache != nil {
		s.cache.DeletePattern(ctx, "products:list:*")
	}

	return nil
}

func (s *productService) GetProduct(ctx context.Context, id generated.IdParam) (*models.Product, error) {
	cacheKey := fmt.Sprintf("product:%s", id)

	// Try to get from cache
	if s.cache != nil {
		var product models.Product
		if err := s.cache.Get(ctx, cacheKey, &product); err == nil {
			return &product, nil
		}
	}

	// Get from database
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Set cache
	if s.cache != nil {
		s.cache.Set(ctx, cacheKey, product, 5*time.Minute)
	}

	return product, nil
}

func (s *productService) ListProducts(ctx context.Context, page, perPage int) ([]models.Product, int64, error) {
	cacheKey := fmt.Sprintf("products:list:%d:%d", page, perPage)

	// Try to get from cache
	if s.cache != nil {
		var result struct {
			Products []models.Product
			Total    int64
		}
		if err := s.cache.Get(ctx, cacheKey, &result); err == nil {
			return result.Products, result.Total, nil
		}
	}

	// Get from database
	products, total, err := s.repo.FindAll(ctx, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	// Set cache
	if s.cache != nil {
		result := struct {
			Products []models.Product
			Total    int64
		}{products, total}
		s.cache.Set(ctx, cacheKey, result, 2*time.Minute)
	}

	return products, total, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id generated.IdParam, product *models.Product) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	product.ID = existing.ID
	product.CreatedAt = existing.CreatedAt

	if err := s.repo.Update(ctx, product); err != nil {
		return err
	}

	// Invalidate cache
	if s.cache != nil {
		s.cache.Delete(ctx, fmt.Sprintf("product:%s", id))
		s.cache.DeletePattern(ctx, "products:list:*")
	}

	return nil
}

func (s *productService) DeleteProduct(ctx context.Context, id generated.IdParam) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	if s.cache != nil {
		s.cache.Delete(ctx, fmt.Sprintf("product:%d", id))
		s.cache.DeletePattern(ctx, "products:list:*")
	}

	return nil
}
