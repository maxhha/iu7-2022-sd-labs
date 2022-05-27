package handlers

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func NewPlaygroundHandler(graphqlEndpoint string) gin.HandlerFunc {
	h := playground.Handler("GraphQL", graphqlEndpoint)

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
