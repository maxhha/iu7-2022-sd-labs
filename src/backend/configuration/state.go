package configuration

import (
	"time"
)

type GraphQLHandlerConfig struct {
	QueryCache                   int // 1000
	AutomaticPersistedQueryCache int // 100
	Introspection                bool
}

func (c *GraphQLHandlerConfig) Equal(other *GraphQLHandlerConfig) bool {
	return c.QueryCache == other.QueryCache &&
		c.AutomaticPersistedQueryCache == other.AutomaticPersistedQueryCache &&
		c.Introspection == other.Introspection
}

type GORMRepositoryConfig struct {
	Database string
	DSN      string
}

func (c *GORMRepositoryConfig) Equal(other *GORMRepositoryConfig) bool {
	return c.Database == other.Database && c.DSN == other.DSN
}

type ConfigurationState struct {
	allowOrigins         []string
	shutdownTimeout      time.Duration
	addr                 string
	graphqlHandlerConfig GraphQLHandlerConfig
	gormRepositoryConfig GORMRepositoryConfig
	signingKey           string
}

type ConfigurationStateBuilder struct {
	state ConfigurationState
}

func NewConfigurationStateBuilder() ConfigurationStateBuilder {
	return ConfigurationStateBuilder{}
}

func (b *ConfigurationStateBuilder) Result() ConfigurationState {
	return b.state
}

func (s *ConfigurationState) AllowOrigins() []string {
	return s.allowOrigins
}

func (b *ConfigurationStateBuilder) SetAllowOrigins(allowOrigins []string) {
	b.state.allowOrigins = allowOrigins
}

func (s *ConfigurationState) ShutdownTimeout() time.Duration {
	return s.shutdownTimeout
}

func (b *ConfigurationStateBuilder) SetShutdownTimeout(shutdownTimeout time.Duration) {
	b.state.shutdownTimeout = shutdownTimeout
}

func (s *ConfigurationState) Addr() string {
	return s.addr
}

func (b *ConfigurationStateBuilder) SetAddr(addr string) {
	b.state.addr = addr
}

func (s *ConfigurationState) GraphQLHandlerConfig() GraphQLHandlerConfig {
	return s.graphqlHandlerConfig
}

func (b *ConfigurationStateBuilder) SetGraphQLHandlerConfig(config GraphQLHandlerConfig) {
	b.state.graphqlHandlerConfig = config
}

func (s *ConfigurationState) GORMRepositoryConfig() GORMRepositoryConfig {
	return s.gormRepositoryConfig
}

func (b *ConfigurationStateBuilder) SetGORMRepositoryConfig(config GORMRepositoryConfig) {
	b.state.gormRepositoryConfig = config
}

func (s *ConfigurationState) SigningKey() string {
	return s.signingKey
}

func (b *ConfigurationStateBuilder) SetSigningKey(signingKey string) {
	b.state.signingKey = signingKey
}
