package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
	"time"
)

func GraphQLLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		var gqlRequest struct {
			OperationName string                 `json:"operationName"`
			Query         string                 `json:"query"`
			Variables     map[string]interface{} `json:"variables"`
		}
		if err := c.BindJSON(&gqlRequest); err == nil {
			// Format the query, operation, and variables as a single JSON string
			formattedQuery := formatGraphQLQuery(gqlRequest.OperationName, gqlRequest.Query, gqlRequest.Variables)
			logger.Info("GraphQL Request",
				zap.String("formattedQuery", formattedQuery),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("duration", time.Since(start)),
			)
		}
	}
}

// formatGraphQLQuery formats the operation, query, and variables into a single JSON string
func formatGraphQLQuery(operation, query string, variables map[string]interface{}) string {
	request := struct {
		OperationName string                 `json:"operationName,omitempty"`
		Query         string                 `json:"query"`
		Variables     map[string]interface{} `json:"variables,omitempty"`
	}{
		OperationName: operation,
		Query:         query,
		Variables:     variables,
	}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return "Error formatting GraphQL query"
	}
	return string(jsonBytes)
}
