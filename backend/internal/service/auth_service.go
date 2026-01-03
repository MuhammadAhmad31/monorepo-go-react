package service

import (
	"backend/internal/generated"
	"backend/internal/models"
	"backend/internal/repository"
	jwt "backend/pkg"
	"context"
	"errors"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, req *generated.RegisterRequest) (*generated.AuthResponse, error)
	Login(ctx context.Context, req *generated.LoginRequest) (*generated.AuthResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(ctx context.Context, req *generated.RegisterRequest) (*generated.AuthResponse, error) {
	// Check if email exists
	existing, err := s.userRepo.FindByEmail(ctx, string(req.Email))
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Email:    string(req.Email),
		Role:     "user",
		IsActive: true,
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := jwt.GenerateToken(user.ID.String(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &generated.AuthResponse{
		Token: token,
		User: generated.UserData{
			Id:       user.ID,
			Name:     user.Name,
			Email:    openapi_types.Email(user.Email),
			Role:     generated.UserDataRole(user.Role),
			IsActive: user.IsActive,
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, req *generated.LoginRequest) (*generated.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, string(req.Email))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	if !user.CheckPassword(string(req.Password)) {
		return nil, errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	token, err := jwt.GenerateToken(user.ID.String(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &generated.AuthResponse{
		Token: token,
		User: generated.UserData{
			Id:       user.ID,
			Name:     user.Name,
			Email:    openapi_types.Email(user.Email),
			Role:     generated.UserDataRole(user.Role),
			IsActive: user.IsActive,
		},
	}, nil
}
