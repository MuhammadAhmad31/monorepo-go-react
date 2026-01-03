package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

// GetOAPISecurityScopes returns security scopes for current request
// This would be populated by oapi-codegen
func GetOAPISecurityScopes(ctx context.Context) []string {
	if gc, ok := ctx.(*gin.Context); ok {
		if scopes, exists := gc.Get("oapi:scopes"); exists {
			return scopes.([]string)
		}
	}
	return nil
}

// SetOAPISecurityScopes sets security scopes for current request
func SetOAPISecurityScopes(c *gin.Context, scopes []string) {
	c.Set("oapi:scopes", scopes)
}
