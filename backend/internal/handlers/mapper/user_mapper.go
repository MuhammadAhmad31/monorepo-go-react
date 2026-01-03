package mapper

import (
	"backend/internal/generated"
	"backend/internal/models"

	types_generated "github.com/oapi-codegen/runtime/types"
)

func ToGeneratedUser(user *models.User) generated.User {
	role := generated.UserRole(user.Role)

	return generated.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     types_generated.Email(user.Email),
		Role:      &role,
		IsActive:  &user.IsActive,
		CreatedAt: &user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
	}
}

func ToGeneratedUsers(users []models.User) []generated.User {
	result := make([]generated.User, len(users))
	for i := range users {
		result[i] = ToGeneratedUser(&users[i])
	}
	return result
}
