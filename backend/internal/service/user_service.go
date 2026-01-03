package service

import (
	"backend/internal/cache"
	"backend/internal/generated"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id generated.IdParam) (*models.User, error)
	ListUsers(ctx context.Context, page, perPage int) ([]models.User, int64, error)
	UpdateUser(ctx context.Context, id generated.IdParam, user *models.User) error
	DeleteUser(ctx context.Context, id generated.IdParam) error
}

type userService struct {
	repo  repository.UserRepository
	cache *cache.RedisCache
}

func NewUserService(repo repository.UserRepository, cache *cache.RedisCache) UserService {
	return &userService{
		repo:  repo,
		cache: cache,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
	// Check if email already exists
	existing, err := s.repo.FindByEmail(ctx, user.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if existing != nil {
		return fmt.Errorf("email already exists")
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	// Invalidate users list cache
	if s.cache != nil {
		s.cache.DeletePattern(ctx, "users:list:*")
	}

	return nil
}

func (s *userService) GetUser(ctx context.Context, id generated.IdParam) (*models.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	// Try to get from cache
	if s.cache != nil {
		var user models.User
		if err := s.cache.Get(ctx, cacheKey, &user); err == nil {
			return &user, nil
		}
	}

	// Get from database
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Set cache
	if s.cache != nil {
		s.cache.Set(ctx, cacheKey, user, 5*time.Minute)
	}

	return user, nil
}

func (s *userService) ListUsers(ctx context.Context, page, perPage int) ([]models.User, int64, error) {
	cacheKey := fmt.Sprintf("users:list:%d:%d", page, perPage)

	// Try to get from cache
	if s.cache != nil {
		var result struct {
			Users []models.User
			Total int64
		}
		if err := s.cache.Get(ctx, cacheKey, &result); err == nil {
			return result.Users, result.Total, nil
		}
	}

	// Get from database
	users, total, err := s.repo.FindAll(ctx, page, perPage)
	if err != nil {
		return nil, 0, err
	}

	// Set cache
	if s.cache != nil {
		result := struct {
			Users []models.User
			Total int64
		}{users, total}
		s.cache.Set(ctx, cacheKey, result, 2*time.Minute)
	}

	return users, total, nil
}

func (s *userService) UpdateUser(ctx context.Context, id generated.IdParam, user *models.User) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	user.ID = existing.ID
	user.CreatedAt = existing.CreatedAt

	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	// Invalidate cache
	if s.cache != nil {
		s.cache.Delete(ctx, fmt.Sprintf("user:%d", id))
		s.cache.DeletePattern(ctx, "users:list:*")
	}

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id generated.IdParam) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	if s.cache != nil {
		s.cache.Delete(ctx, fmt.Sprintf("user:%d", id))
		s.cache.DeletePattern(ctx, "users:list:*")
	}

	return nil
}
