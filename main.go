package main

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/blnto/blnto_service/internal"
	"github.com/blnto/blnto_service/internal/infrastructure/graphql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Defining the Graphql handler
func graphqlHandler(app *internal.App) gin.HandlerFunc {
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: app.Resolver}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func setupRoutes(router *gin.Engine, app *internal.App) {
	router.GET("/", playgroundHandler())
	router.POST("/query", GraphQLLogger(app.Logger), graphqlHandler(app))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	app, err := internal.InitializeDependencies()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	// artistApi.StartTokenRefreshScheduler(app.DB)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:52322"},
		AllowMethods:     []string{"POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Setup routes
	setupRoutes(router, app)
	logFile := app.Loggerfile
	logger := app.Logger

	defer logFile.Close()
	defer logger.Sync()
	// Start the application
	if err := router.Run(":8080"); err != nil {
		// handle error, perhaps log it and/or gracefully shut down
		log.Fatalf("Failed to run server: %v", err)
	}
}

func GraphQLLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read and store the request body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("Error reading request body", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error reading request body"})
			return
		}

		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// Unmarshal the body to the gqlRequest struct for logging
		var gqlRequest struct {
			OperationName string                 `json:"operationName"`
			Query         string                 `json:"query"`
			Variables     map[string]interface{} `json:"variables"`
		}
		if err := json.Unmarshal(bodyBytes, &gqlRequest); err != nil {
			logger.Error("Error unmarshalling request body", zap.Error(err))
		}

		// Log the GraphQL request
		logger.Info("GraphQL Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("operation", gqlRequest.OperationName),
			zap.String("query", gqlRequest.Query),
			zap.Any("variables", gqlRequest.Variables),
		)

		start := time.Now()
		c.Next()

		// Log the response status and duration after handling the request
		duration := time.Since(start)
		logger.Info("GraphQL Response",
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
		)
	}
}
