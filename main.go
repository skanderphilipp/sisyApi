package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/skanderphilipp/sisyApi/internal"
	"github.com/skanderphilipp/sisyApi/internal/infrastructure/graphql"
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
	router.POST("/query", graphqlHandler(app))
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func main() {
	// Initialize the application with Wire
	app, err := internal.InitializeDependencies()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	router := gin.Default()
	// Setup routes
	setupRoutes(router, app)

	// Start the application
	if err := router.Run(":8080"); err != nil {
		// handle error, perhaps log it and/or gracefully shut down
		log.Fatalf("Failed to run server: %v", err)
	}
}
