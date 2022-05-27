package handlers

import (
	"iu7-2022-sd-labs/configuration"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func NewGraphQLHandler(config configuration.GraphQLHandlerConfig, es graphql.ExecutableSchema) (gin.HandlerFunc, error) {
	h := handler.New(es)

	h.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})
	h.AddTransport(transport.MultipartForm{})

	if config.QueryCache > 0 {
		h.SetQueryCache(lru.New(config.QueryCache))
	}

	if config.Introspection {
		h.Use(extension.Introspection{})
	}

	if config.AutomaticPersistedQueryCache > 0 {
		h.Use(extension.AutomaticPersistedQuery{
			Cache: lru.New(config.AutomaticPersistedQueryCache),
		})
	}

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}, nil
}
