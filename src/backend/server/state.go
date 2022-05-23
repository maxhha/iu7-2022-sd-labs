package server

import (
	"fmt"
	"iu7-2022-sd-labs/adapters/consumer_validator"
	"iu7-2022-sd-labs/adapters/dataloader"
	"iu7-2022-sd-labs/adapters/event_bus"
	"iu7-2022-sd-labs/adapters/gorm_repositories"
	"iu7-2022-sd-labs/adapters/jwt_auth"
	"iu7-2022-sd-labs/adapters/offer_pay_service"
	"iu7-2022-sd-labs/app"
	"iu7-2022-sd-labs/buisness/interactors"
	"iu7-2022-sd-labs/buisness/ports/services"
	"iu7-2022-sd-labs/configuration"
	"iu7-2022-sd-labs/server/generated"
	"iu7-2022-sd-labs/server/handlers"
	"iu7-2022-sd-labs/server/ports"
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

	signingKey                   string
	repoConfig                   configuration.GORMRepositoryConfig
	repo                         *gorm_repositories.GORMRepository
	dataloader                   *dataloader.DataLoader
	auth                         *jwt_auth.JWTAuth
	consumerFormValidatorService services.ConsumerFormValidatorService
	eventBus                     *event_bus.EventBus
	offerPayService              *offer_pay_service.OfferPayService
	auctionInteractor            *interactors.AuctionInteractor
	bidStepTableInteractor       *interactors.BidStepTableInteractor
	consumerInteractor           *interactors.ConsumerInteractor
	offerInteractor              *interactors.OfferInteractor
	organizerInteractor          *interactors.OrganizerInteractor
	productInteractor            *interactors.ProductInteractor
	roomInteractor               *interactors.RoomInteractor
	resolver                     generated.ResolverRoot
	schema                       graphql.ExecutableSchema
	graphqlConfig                configuration.GraphQLHandlerConfig
	graphqlHandeler              gin.HandlerFunc
	engine                       *gin.Engine
	server                       *graceful.Server
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
	eventBus := event_bus.NewEventBus()
	consumerFormValidator := consumer_validator.NewConsumerFormValidatorService()
	offerPayService := offer_pay_service.NewOfferPayService()

	s := &ServerState{
		logger:                       logger,
		runResultChan:                errorChan,
		serverErrorChan:              errorChan,
		eventBus:                     &eventBus,
		consumerFormValidatorService: &consumerFormValidator,
		offerPayService:              &offerPayService,
	}

	var err error
	s, err = s.update(&config)
	if err != nil {
		return nil, Wrap(err, "update")
	}

	return s, nil
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

func (s *ServerState) Update(config *configuration.ConfigurationState) (app.AppState, error) {
	return s.update(config)
}

func (s *ServerState) update(config *configuration.ConfigurationState) (*ServerState, error) {
	s.logger.Println("state update")
	nextState := *s

	callables := []struct {
		update func(prev *ServerState, config *configuration.ConfigurationState) error
		name   string
	}{
		{nextState.updateRepo, "Repo"},
		{nextState.updateDataLoader, "DataLoader"},
		{nextState.updateAuth, "Auth"},
		{nextState.updateAuctionInteractor, "AuctionInteractor"},
		{nextState.updateBidStepTableInteractor, "BidStepTableInteractor"},
		{nextState.updateConsumerInteractor, "ConsumerInteractor"},
		{nextState.updateOfferInteractor, "OfferInteractor"},
		{nextState.updateOrganizerInteractor, "OrganizerInteractor"},
		{nextState.updateProductInteractor, "ProductInteractor"},
		{nextState.updateRoomInteractor, "RoomInteractor"},
		{nextState.updateResolver, "Resolver"},
		{nextState.updateSchema, "Schema"},
		{nextState.updateGraphqlHandler, "GraphqlHandler"},
		{nextState.updateEngine, "Engine"},
		{nextState.updateServer, "Server"},
	}
	for _, v := range callables {
		if err := v.update(s, config); err != nil {
			return nil, Wrapf(err, "update %s", v.name)
		}
	}

	return &nextState, nil
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
		func() ports.Auth {
			return s.auth
		},
		s.logger,
		func() ports.DataLoader {
			return s.dataloader
		},
	)

	return nil
}

func newEngine(
	config *configuration.ConfigurationState,
	graphqlHandler gin.HandlerFunc,
	playgroundHandler gin.HandlerFunc,
	authGetter func() ports.Auth,
	logger *log.Logger,
	dataLoaderGetter func() ports.DataLoader,
) *gin.Engine {
	allowOrigins := config.AllowOrigins()
	engine := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = allowOrigins
	corsConfig.AllowMethods = []string{"POST, GET, OPTIONS"}
	engine.Use(cors.New(corsConfig))

	engine.Use(newAuthMiddleware(authGetter, logger))
	engine.Use(newDataLoaderMiddleware(dataLoaderGetter))

	engine.Any(graphqlRoute, graphqlHandler)
	engine.GET(graphiqlRoute, playgroundHandler)

	return engine
}

func newAuthMiddleware(auth func() ports.Auth, logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()
		a := auth()

		tokens, ok := c.Request.Header["Authorization"]

		if !ok || len(tokens) == 0 {
			return
		}

		token := tokens[0]

		consumer, err := a.ParseConsumerToken(token)
		if err == nil {
			ctx := ports.WithConsumer(c.Request.Context(), consumer)
			c.Request = c.Request.WithContext(ctx)
		} else if !Is(err, ports.ErrWrongSubject) {
			logger.Printf("error on parse consuemr token: %v\n", err)
			return
		}

		organizer, err := a.ParseOrganizerToken(token)
		if err == nil {
			ctx := ports.WithOrganizer(c.Request.Context(), organizer)
			c.Request = c.Request.WithContext(ctx)
		} else if !Is(err, ports.ErrWrongSubject) {
			logger.Printf("error on parse organizer token: %v\n", err)
			return
		}
	}
}

func newDataLoaderMiddleware(getDataloader func() ports.DataLoader) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()
		dataloader := getDataloader()
		ctx := dataloader.WithNewLoader(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)
	}
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

func (s *ServerState) updateDataLoader(prev *ServerState, config *configuration.ConfigurationState) error {
	repoChanged := prev.repo != s.repo
	shouldReset := repoChanged

	if !shouldReset {
		return nil
	}

	s.dataloader = dataloader.NewDataLoader(s.repo)
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

func (s *ServerState) updateConsumerInteractor(prev *ServerState, config *configuration.ConfigurationState) error {
	repoChanged := s.repo != prev.repo
	eventBusChanged := s.eventBus != prev.eventBus
	shouldReset := repoChanged || eventBusChanged

	if !shouldReset {
		return nil
	}

	consumerInteractor := interactors.NewConsumerInteractor(s.repo, s.eventBus, s.consumerFormValidatorService)
	s.consumerInteractor = &consumerInteractor
	return nil
}

func (s *ServerState) updateOrganizerInteractor(prev *ServerState, config *configuration.ConfigurationState) error {
	repoChanged := s.repo != prev.repo
	shouldReset := repoChanged

	if !shouldReset {
		return nil
	}

	organizerInteractor := interactors.NewOrganizerInteractor(s.repo.Organizer())
	s.organizerInteractor = &organizerInteractor
	return nil
}

func (s *ServerState) updateRoomInteractor(prev *ServerState, config *configuration.ConfigurationState) error {
	repoChanged := s.repo != prev.repo
	shouldReset := repoChanged

	if !shouldReset {
		return nil
	}

	roomInteractor := interactors.NewRoomInteractor(s.repo.Organizer(), s.repo.Room())
	s.roomInteractor = &roomInteractor
	return nil
}

func (s *ServerState) updateProductInteractor(prev *ServerState, config *configuration.ConfigurationState) error {
	repoChanged := s.repo != prev.repo
	shouldReset := repoChanged

	if !shouldReset {
		return nil
	}

	productInteractor := interactors.NewProductInteractor(s.repo.Organizer(), s.repo.Product())
	s.productInteractor = &productInteractor
	return nil
}

func (s *ServerState) updateAuctionInteractor(prev *ServerState, config *configuration.ConfigurationState) error {
	repoChanged := s.repo != prev.repo
	eventBusChanged := s.eventBus != prev.eventBus
	shouldReset := repoChanged || eventBusChanged

	if !shouldReset {
		return nil
	}

	auctionInteractor := interactors.NewAuctionInteractor(s.repo, s.eventBus)
	s.auctionInteractor = &auctionInteractor
	return nil
}

func (s *ServerState) updateBidStepTableInteractor(prev *ServerState, config *configuration.ConfigurationState) error {
	repoChanged := s.repo != prev.repo
	shouldReset := repoChanged

	if !shouldReset {
		return nil
	}

	bidStepTableInteractor := interactors.NewBidStepTableInteractor(s.repo)
	s.bidStepTableInteractor = &bidStepTableInteractor
	return nil
}

func (s *ServerState) updateOfferInteractor(prev *ServerState, config *configuration.ConfigurationState) error {
	repoChanged := s.repo != prev.repo
	eventBusChanged := s.eventBus != prev.eventBus
	offerPayServiceChanged := s.offerPayService != prev.offerPayService

	shouldReset := repoChanged || eventBusChanged || offerPayServiceChanged

	if !shouldReset {
		return nil
	}

	offerInteractor := interactors.NewOfferInteractor(s.repo, s.eventBus, s.offerPayService)
	s.offerInteractor = &offerInteractor
	return nil
}

func (s *ServerState) updateResolver(prev *ServerState, config *configuration.ConfigurationState) error {
	organizerInteractorChanged := s.organizerInteractor != prev.organizerInteractor
	consumerInteractorChanged := s.consumerInteractor != prev.consumerInteractor
	roomInteractorChanged := s.roomInteractor != prev.roomInteractor
	dataLoaderChanged := s.dataloader != prev.dataloader
	authChanged := s.auth != prev.auth
	eventBusChanged := s.eventBus != prev.eventBus
	shouldReset := authChanged ||
		organizerInteractorChanged ||
		consumerInteractorChanged ||
		roomInteractorChanged ||
		dataLoaderChanged ||
		eventBusChanged

	if !shouldReset {
		return nil
	}

	resolver := resolvers.New(
		s.auctionInteractor,
		s.bidStepTableInteractor,
		s.consumerInteractor,
		s.offerInteractor,
		s.organizerInteractor,
		s.productInteractor,
		s.roomInteractor,
		s.auth,
		s.dataloader,
		s.eventBus,
		s.logger,
	)
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
