package middleware

import (
	"github.com/gin-gonic/gin"
)

/*
	Used for common connections to the database
*/

// DatabaseMiddleware is used for connections to the database
func DatabaseMiddleware(dbPass, dbURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db_pass", dbPass)
		c.Set("db_url", dbURL)
		// only call this inside middleware
		// it's purpose is to execute any pending handlers
		c.Next()
	}
}