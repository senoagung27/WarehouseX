package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Role hierarchy: admin > supervisor > staff, auditor is read-only
var roleHierarchy = map[string]int{
	"admin":      4,
	"supervisor": 3,
	"staff":      2,
	"auditor":    1,
}

// RequireRole checks if the user has the minimum required role
func RequireRole(minRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := GetUserRole(c)

		userLevel, ok := roleHierarchy[userRole]
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid role"})
			c.Abort()
			return
		}

		requiredLevel, ok := roleHierarchy[minRole]
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid role configuration"})
			c.Abort()
			return
		}

		if userLevel < requiredLevel {
			c.JSON(http.StatusForbidden, gin.H{
				"error":         "insufficient permissions",
				"required_role": minRole,
				"your_role":     userRole,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRoles checks if the user has one of the specified roles
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := GetUserRole(c)

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error":          "insufficient permissions",
			"required_roles": roles,
			"your_role":      userRole,
		})
		c.Abort()
	}
}
