// rbac.go
package middleware

import (
	"backend/internal/generated"
	jwt "backend/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// OpenAPISecurityMiddleware enforces security rules from OpenAPI spec
// Uses auto-generated RouteSecurity map from contracts/openapi.yaml
func OpenAPISecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		// Get security info from generated map
		secInfo := getRouteSecurityInfo(path, method)

		// Public endpoint - skip auth
		if secInfo.IsPublic {
			c.Next()
			return
		}

		// Protected endpoint - validate JWT
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, generated.Error{
				Message: "authorization required",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, generated.Error{
				Message: "invalid authorization format",
			})
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, generated.Error{
				Message: "invalid token",
			})
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		// Empty scopes = any authenticated user is allowed
		if len(secInfo.RequiredScopes) == 0 {
			c.Next()
			return
		}

		// Check if user role matches required scopes
		for _, scope := range secInfo.RequiredScopes {
			if claims.Role == scope {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, generated.Error{
			Message: "insufficient permissions",
		})
	}
}

// getRouteSecurityInfo retrieves security info from generated.RouteSecurity
func getRouteSecurityInfo(path, method string) generated.RouteSecurityInfo {
	// Check exact path match
	if methods, ok := generated.RouteSecurity[path]; ok {
		if secInfo, ok := methods[method]; ok {
			return secInfo
		}
	}

	// Check path with dynamic parameter (e.g., /api/v1/users/123)
	for routePath, methods := range generated.RouteSecurity {
		if matchDynamicRoute(path, routePath) {
			if secInfo, ok := methods[method]; ok {
				return secInfo
			}
		}
	}

	// Default: require authentication but no specific role
	return generated.RouteSecurityInfo{
		IsPublic:       false,
		RequiredScopes: []string{},
	}
}

// matchDynamicRoute checks if actual path matches route pattern
// Example: /api/v1/users/123 matches /api/v1/users/{id}
func matchDynamicRoute(actualPath, routePattern string) bool {
	if !strings.Contains(routePattern, "{") {
		return false
	}

	actualParts := strings.Split(actualPath, "/")
	patternParts := strings.Split(routePattern, "/")

	if len(actualParts) != len(patternParts) {
		return false
	}

	for i, patternPart := range patternParts {
		// {id}, {productId}, etc. match any value
		if strings.HasPrefix(patternPart, "{") && strings.HasSuffix(patternPart, "}") {
			continue
		}

		if actualParts[i] != patternPart {
			return false
		}
	}

	return true
}
