package server

import (
	"fmt"
	"iu7-2022-sd-labs/adapters/gorm_repositories"
	"iu7-2022-sd-labs/adapters/jwt_auth"
	"iu7-2022-sd-labs/app"
	"iu7-2022-sd-labs/buisness/interactors"
	"iu7-2022-sd-labs/configuration"
	"iu7-2022-sd-labs/server/generated"
	"iu7-2022-sd-labs/server/handlers"
	"iu7-2022-sd-labs/server/resolvers"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/tylerb/graceful.v1"
)

const graphqlRoute = "/graphql"
const graphiqlRoute = "/graphiql"

type ServerState struct {
	logger          *log.Logger
	stopped         bool
	runResultChan   chan error
	serverErrorChan chan error

	allowOrigins []string

	signingKey          string
	repoConfig          configuration.GORMRepositoryConfig
	repo                *gorm_repositories.GORMRepository
	auth                *jwt_auth.JWTAuth
	organizerInteractor *interactors.OrganizerInteractor
	resolver            *resolvers.Resolver
	schema              graphql.ExecutableSchema
	graphqlConfig       configuration.GraphQLHandlerConfig
	graphqlHandeler     gin.HandlerFunc
	engine              *gin.Engine
	server              *graceful.Server
}

func isStringArrayEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, ai := range a {
		if b[i] != ai {
			return false
		}
	}

	return true
}

func NewServerState(configuration *configuration.Configuration) (*ServerState, error) {
	logger := configuration.Logger()
	logger.Println("state init")
	config := configuration.Current()

	errorChan := make(chan error)

	s := &ServerState{
		logger:          logger,
		runResultChan:   errorChan,
		serverErrorChan: errorChan,
	}

	var err error
	s, err = s.update(&config)
	if err != nil {
		return nil, Wrap(err, "update")
	}

	return s, nil
}

func (s *ServerState) Update(config *configuration.ConfigurationState) (app.AppState, error) {
	return s.update(config)
}

func (s *ServerState) update(config *configuration.ConfigurationState) (*ServerState, error) {
	s.logger.Println("state update")
	nextState := *s

	if err := nextState.updateRepo(s, config); err != nil {
		return nil, Wrap(err, "update repo")
	}
	if err := nextState.updateAuth(s, config); err != nil {
		return nil, Wrap(err, "update auth")
	}
	if err := nextState.updateOrganizerInteractor(s, config); err != nil {
		return nil, Wrap(err, "update organizer interactor")
	}
	if err := nextState.updateResolver(s, config); err != nil {
		return nil, Wrap(err, "update resolver")
	}
	if err := nextState.updateSchema(s, config); err != nil {
		return nil, Wrap(err, "update schema")
	}
	if err := nextState.updateGraphqlHandler(s, config); err != nil {
		return nil, Wrap(err, "update graphql handler")
	}
	if err := nextState.updateEngine(s, config); err != nil {
		return nil, Wrap(err, "update engine")
	}
	if err := nextState.updateServer(s, config); err != nil {
		return nil, Wrap(err, "update server")
	}

	return &nextState, nil
}

func (s *ServerState) ChangeTo(state app.AppState) error {
	s.logger.Println("change state")
	nextState, ok := state.(*ServerState)
	if !ok {
		return fmt.Errorf("cant change to %#v", nextState)
	}

	if s.server != nextState.server {
		s.serverErrorChan = make(chan error)
		go s.Stop()
		if err := <-s.serverErrorChan; err != nil {
			s.logger.Printf("state change to: prev server stop: %v\n", err)
		}
		go nextState.listenAndServe()
	}

	*s = *nextState

	return nil
}

func (s *ServerState) Stop() {
	if s.stopped {
		return
	}

	s.logger.Println("state stop")
	s.stopped = true
	s.server.Stop(s.server.Timeout)
}

func (s *ServerState) Run() error {
	go s.listenAndServe()
	return <-s.runResultChan
}

func (s *ServerState) listenAndServe() {
	s.serverErrorChan <- s.server.ListenAndServe()
}

func (s *ServerState) updateEngine(prev *ServerState, config *configuration.ConfigurationState) error {
	s.allowOrigins = config.AllowOrigins()
	isAllowOriginsEmpty := len(s.allowOrigins) == 0
	isAllowOriginsChanged := !isStringArrayEqual(s.allowOrigins, prev.allowOrigins)

	shouldReset := isAllowOriginsEmpty || isAllowOriginsChanged

	if !shouldReset {
		return nil
	}

	s.engine = newEngine(
		config,
		func(ctx *gin.Context) {
			s.graphqlHandeler(ctx)
		},
		handlers.NewPlaygroundHandler(graphqlRoute),
	)

	return nil
}

func newEngine(
	config *configuration.ConfigurationState,
	graphqlHandler gin.HandlerFunc,
	playgroundHandler gin.HandlerFunc,
) *gin.Engine {
	allowOrigins := config.AllowOrigins()
	engine := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = allowOrigins
	corsConfig.AllowMethods = []string{"POST, GET, OPTIONS"}
	engine.Use(cors.New(corsConfig))

	// engine.Use(auth.New(&db))

	engine.Any(graphqlRoute, graphqlHandler)
	engine.GET(graphiqlRoute, playgroundHandler)

	return engine
}

func (s *ServerState) updateServer(prev *ServerState, config *configuration.ConfigurationState) error {
	addrChanged := prev.server == nil || config.Addr() != prev.server.Addr
	engineChanged := s.engine != prev.engine
	shouldReset := engineChanged || addrChanged

	if shouldReset {
		s.server = newServer(
			config,
			s.logger,
			s.engine,
		)
		return nil
	}

	s.server.Timeout = config.ShutdownTimeout()
	return nil
}

func newServer(
	config *configuration.ConfigurationState,
	logger *log.Logger,
	engine *gin.Engine,
) *graceful.Server {
	return &graceful.Server{
		Timeout:          config.ShutdownTimeout(),
		NoSignalHandling: true,
		Logger:           logger,

		Server: &http.Server{
			Addr:    config.Addr(),
			Handler: engine,
		},
	}
}

func (s *ServerState) updateRepo(prev *ServerState, config *configuration.ConfigurationState) error {
	s.repoConfig = config.GORMRepositoryConfig()

	configChanged := !prev.repoConfig.Equal(&s.repoConfig)
	shouldReset := configChanged

	if !shouldReset {
		return nil
	}
	repo, err := gorm_repositories.NewFromConfig(s.repoConfig)
	if err != nil {
		return Wrap(err, "new gorm repository from config")
	}

	s.repo = &repo
	return nil
}

func (s *ServerState) updateAuth(prev *ServerState, config *configuration.ConfigurationState) error {
	s.signingKey = config.SigningKey()

	keyChanged := s.signingKey != prev.signingKey
	repoChanged := s.repo != prev.repo
	shouldReset := keyChanged || repoChanged

	if !shouldReset {
		return nil
	}

	auth := jwt_auth.NewJWTAuth(
		s.signingKey,
		s.repo.Organizer(),
		s.repo.Consumer(),
	)

	s.auth = &auth
	return nil
}

func (s *ServerState) updateOrganizerInteractor(prev *ServerState, config *configuration.ConfigurationState) error {
	repoChanged := s.repo != prev.repo
	shouldReset := repoChanged

	if !shouldReset {
		return nil
	}

	orgainzerInteractor := interactors.NewOrganizerInteractor(s.repo.Organizer())
	s.organizerInteractor = &orgainzerInteractor
	return nil
}

func (s *ServerState) updateResolver(prev *ServerState, config *configuration.ConfigurationState) error {
	organizerInteractorChanged := s.organizerInteractor != prev.organizerInteractor
	authChanged := s.auth != prev.auth
	shouldReset := organizerInteractorChanged || authChanged

	if !shouldReset {
		return nil
	}

	resolver := resolvers.New(s.organizerInteractor, s.auth)
	s.resolver = &resolver

	return nil
}

func (s *ServerState) updateSchema(prev *ServerState, config *configuration.ConfigurationState) error {
	resolverChanged := s.resolver != prev.resolver
	shouldReset := resolverChanged

	if !shouldReset {
		return nil
	}

	schemaConfig := generated.Config{Resolvers: s.resolver}
	// config.Directives.HasRole = hasRoleDirective(roleCheckers)
	s.schema = generated.NewExecutableSchema(schemaConfig)

	return nil
}

func (s *ServerState) updateGraphqlHandler(prev *ServerState, config *configuration.ConfigurationState) error {
	s.graphqlConfig = config.GraphQLHandlerConfig()
	graphqlConfigChanged := !s.graphqlConfig.Equal(&prev.graphqlConfig)
	schemaChanged := s.schema != prev.schema

	shouldReset := graphqlConfigChanged || schemaChanged

	if !shouldReset {
		return nil
	}

	handler, err := handlers.NewGraphQLHandler(s.graphqlConfig, s.schema)
	if err != nil {
		return Wrap(err, "new graphql handler")
	}

	s.graphqlHandeler = handler
	return nil
}
