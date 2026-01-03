package handlers

import (
	"backend/internal/generated"
	"backend/internal/handlers/mapper"
	"backend/internal/models"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) ListUsers(c *gin.Context, params generated.ListUsersParams) {
	page := 1
	perPage := 10

	if params.Page != nil {
		page = *params.Page
	}
	if params.PerPage != nil {
		perPage = *params.PerPage
	}

	users, total, err := h.service.ListUsers(c.Request.Context(), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Failed to fetch users",
		})
		return
	}

	totalInt := int(total)

	c.JSON(http.StatusOK, gin.H{
		"data": mapper.ToGeneratedUsers(users),
		"meta": generated.Meta{
			Page:    &page,
			PerPage: &perPage,
			Total:   &totalInt,
		},
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req generated.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, generated.Error{
			Message: "Invalid request body",
		})
		return
	}

	user := &models.User{
		Name:     req.Name,
		Email:    string(req.Email),
		Role:     "user",
		IsActive: true,
	}

	if err := h.service.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, generated.Error{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": mapper.ToGeneratedUser(user),
	})
}

func (h *UserHandler) GetUser(c *gin.Context, id generated.IdParam) {
	user, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, generated.Error{
				Message: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Failed to fetch user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": mapper.ToGeneratedUser(user),
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context, id generated.IdParam) {
	var req generated.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, generated.Error{
			Message: "Invalid request body",
		})
		return
	}

	// Get existing user first
	existing, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, generated.Error{
				Message: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Failed to fetch user",
		})
		return
	}

	// Update fields
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Email != nil {
		existing.Email = string(*req.Email)
	}
	if req.Role != nil {
		existing.Role = string(*req.Role)
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	if err := h.service.UpdateUser(c.Request.Context(), id, existing); err != nil {
		c.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": mapper.ToGeneratedUser(existing),
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context, id generated.IdParam) {
	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, generated.Error{
				Message: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, generated.Error{
			Message: "Failed to delete user",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
