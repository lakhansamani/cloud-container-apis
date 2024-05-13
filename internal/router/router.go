package router

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	"github.com/lakhansamani/cloud-container/graph"
	"github.com/lakhansamani/cloud-container/internal/middleware"
)

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: graph.NewResolver(),
	}))
	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/v1/graphql")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// HealthCheckHandler returns a 200 OK
func HealthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte("OK"))
	}
}

// Defining the router for the server
func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(middleware.DefaultStructuredLogger())
	router.Use(middleware.GinContextToContextMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.AuthorizationMiddleware())
	router.Use(gin.Recovery())
	router.GET("/health", HealthCheckHandler())
	// version 1
	apiV1 := router.Group("/v1")
	apiV1.GET("/", playgroundHandler())
	apiV1.POST("/graphql", graphqlHandler())
	return router
}
