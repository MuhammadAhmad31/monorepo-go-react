package handlers

import (
	"backend/internal/generated"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req generated.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, generated.Error{
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, generated.Error{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req generated.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, generated.Error{
			Message: err.Error(),
		})
		return
	}

	response, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, generated.Error{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, generated.Error{
			Message: "missing user_id",
		})
		return
	}

	emailVal, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, generated.Error{
			Message: "missing email",
		})
		return
	}

	roleVal, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusUnauthorized, generated.Error{
			Message: "missing role",
		})
		return
	}

	userID := uuid.MustParse(userIDVal.(string))
	userEmail := openapi_types.Email(emailVal.(string))
	userRole := roleVal.(string)

	c.JSON(http.StatusOK, generated.MeResponse{
		UserId: &userID,
		Email:  &userEmail,
		Role:   &userRole,
	})
}
