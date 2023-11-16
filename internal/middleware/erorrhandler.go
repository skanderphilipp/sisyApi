package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process the request

		// Check if there are any errors
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// Log or handle the error
				log.Printf("Error: %v", e.Err)
			}

			// Respond with a generic error message
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		}
	}
}
