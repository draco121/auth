package main

import (
	"authentication/database"
	"authentication/graph"
	"authentication/graph/generated"
	"context"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func setLogger() {
	log.SetFormatter(&log.JSONFormatter{})
}

func setUserContext() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		items := graph.ContextItems{Database: database.Connect()}
		sessionid := ctx.Request.Header.Get("sessionid")
		if sessionid != "" {
			items.Sessionid = &sessionid
		}
		c := context.WithValue(ctx.Request.Context(), "context_items", items)
		ctx.Request = ctx.Request.WithContext(c)
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
	err := godotenv.Load("config.env")
	if err != nil {
		panic("failed to start application, could not load environment configuration.")
	}
	setLogger()
	port := os.Getenv("PORT")
	r := gin.Default()
	gin.SetMode(os.Getenv("GIN_MODE"))
	r.Use(setUserContext())
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())
	r.Run(":" + port)
}
