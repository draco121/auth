package main

import (
	"auth/graph"
	"auth/graph/generated"
	"auth/startup"
	"context"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func setUserContext() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		sessionid := ctx.Request.Header.Get("sessionid")
		if sessionid != "" {
			c := context.WithValue(ctx.Request.Context(), "sessionid", sessionid)
			ctx.Request = ctx.Request.WithContext(c)
		}
		ctx.Next()
	}
}

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

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

func main() {
	startup.Initialize()
	port := startup.Config.Port
	if port == "" {
		port = os.Getenv("PORT")
	}
	r := gin.Default()
	r.Use(setUserContext())
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())
	r.Run(":" + port)
}
