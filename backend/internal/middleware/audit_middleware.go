package middleware

import (
	"finara-backend/internal/models"
	"finara-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func AuditMiddleware(auditService services.AuditLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request first
		c.Next()

		// Skip audit for certain paths
		if shouldSkipAudit(c.FullPath()) {
			return
		}

		// Get user info from context
		companyID, companyExists := c.Get("company_id")
		userID, userExists := c.Get("user_id")

		if !companyExists || !userExists {
			return
		}

		// Determine action based on HTTP method
		action := getActionFromMethod(c.Request.Method)

		// Get module from path
		module := getModuleFromPath(c.FullPath())

		// Get IP and User Agent
		ipAddress := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// Create audit log
		auditLog := &models.AuditLog{
			CompanyID:   companyID.(uint),
			UserID:      userID.(uint),
			Action:      action,
			Module:      module,
			IPAddress:   ipAddress,
			UserAgent:   userAgent,
			Description: generateDescription(c.Request.Method, c.FullPath()),
		}

		// Log the action (async)
		go auditService.LogAction(auditLog)
	}
}

func shouldSkipAudit(path string) bool {
	skipPaths := []string{
		"/api/v1/health",
		"/api/v1/auth/login",
		"/api/v1/auth/register",
		"/api/v1/profile",
	}

	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}

	return false
}

func getActionFromMethod(method string) models.AuditAction {
	switch method {
	case "POST":
		return models.ActionCreate
	case "PUT", "PATCH":
		return models.ActionUpdate
	case "DELETE":
		return models.ActionDelete
	case "GET":
		return models.ActionView
	default:
		return models.ActionView
	}
}

func getModuleFromPath(path string) string {
	// Extract module from path
	// /api/v1/journals -> journals
	// /api/v1/accounts -> accounts
	if len(path) > 8 {
		parts := []rune(path)
		module := ""
		startIndex := 8 // Skip "/api/v1/"
		
		for i := startIndex; i < len(parts); i++ {
			if parts[i] == '/' {
				break
			}
			module += string(parts[i])
		}
		
		return module
	}
	
	return "unknown"
}

func generateDescription(method, path string) string {
	action := ""
	switch method {
	case "POST":
		action = "Created"
	case "PUT", "PATCH":
		action = "Updated"
	case "DELETE":
		action = "Deleted"
	case "GET":
		action = "Viewed"
	}

	module := getModuleFromPath(path)
	return action + " " + module
}