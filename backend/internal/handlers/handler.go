package handlers

import (
	"backend/internal/service"

	"github.com/gin-gonic/gin"
)

type CombinedHandler struct {
	*UserHandler
	*ProductHandler
	*AuthHandler
}

func NewCombinedHandler(
	userService service.UserService,
	productService service.ProductService,
	authService service.AuthService,
) *CombinedHandler {
	return &CombinedHandler{
		UserHandler:    NewUserHandler(userService),
		ProductHandler: NewProductHandler(productService),
		AuthHandler:    NewAuthHandler(authService),
	}
}

// Helper method to get user context from middleware
func GetUserContext(c *gin.Context) (userID, email, role string) {
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(string)
	}
	if em, exists := c.Get("email"); exists {
		email = em.(string)
	}
	if r, exists := c.Get("role"); exists {
		role = r.(string)
	}
	return
}

// Helper method to check if user is authenticated
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// Helper method to check if user has specific role
func HasRole(c *gin.Context, requiredRole string) bool {
	if role, exists := c.Get("role"); exists {
		return role.(string) == requiredRole
	}
	return false
}
