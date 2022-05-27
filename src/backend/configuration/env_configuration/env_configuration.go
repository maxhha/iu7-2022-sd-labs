package env_configuration

import (
	"fmt"
	"iu7-2022-sd-labs/configuration"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type envConfiguration struct {
	CheckRepeatTime                     time.Duration
	AllowOrigins                        []string
	ShutdownTimeout                     time.Duration
	Addr                                string
	GraphQLQueryCache                   int
	GraphQLAutomaticPersistedQueryCache int
	GraphQLIntrospection                bool
	GormDatabase                        string
	GormDSN                             string
	SigningKey                          string
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

func (c *envConfiguration) Equal(other *envConfiguration) bool {
	return c.CheckRepeatTime == other.CheckRepeatTime &&
		isStringArrayEqual(c.AllowOrigins, other.AllowOrigins) &&
		c.ShutdownTimeout == other.ShutdownTimeout &&
		c.Addr == other.Addr &&
		c.GraphQLQueryCache == other.GraphQLQueryCache &&
		c.GraphQLAutomaticPersistedQueryCache == other.GraphQLAutomaticPersistedQueryCache &&
		c.GraphQLIntrospection == other.GraphQLIntrospection &&
		c.GormDatabase == other.GormDatabase &&
		c.GormDSN == other.GormDSN &&
		c.SigningKey == other.SigningKey
}

func NewEnvConfigurationSource(logger *log.Logger) (<-chan configuration.ConfigurationState, func(), error) {
	config, err := readConfiguration()
	if err != nil {
		return nil, nil, fmt.Errorf("read configuration: %w", err)
	}

	ch := make(chan configuration.ConfigurationState, 1)
	done := make(chan struct{})
	ch <- config.state()

	go func() {
		defer close(ch)

		checkRepeatTime := config.CheckRepeatTime
		timer := time.NewTimer(checkRepeatTime)
		<-timer.C

		for {
			time.NewTimer(time.Millisecond)

			select {
			case <-timer.C:
				logger.Printf("read config file\n")
				newConfig, err := readConfiguration()
				if err != nil {
					logger.Printf("read configuration: %v\n", err)
				}

				if newConfig.Equal(&config) {
					logger.Println("configuration unchanged")
					continue
				}

				config = newConfig

				checkRepeatTime = config.CheckRepeatTime
				ch <- config.state()
				timer.Reset(checkRepeatTime)
			case <-done:
				return
			}
		}
	}()

	cancel := func() {
		done <- struct{}{}
	}

	return ch, cancel, nil
}

func readConfiguration() (envConfiguration, error) {
	var err error
	config := envConfiguration{}

	s, ok := os.LookupEnv("CHECK_REPEAT_TIME")
	if !ok {
		return config, fmt.Errorf("CHECK_REPEAT_TIME not present in environment")
	}

	config.CheckRepeatTime, err = time.ParseDuration(s)
	if err != nil {
		return config, fmt.Errorf("fail convert CHECK_REPEAT_TIME: %w", err)
	}

	s, ok = os.LookupEnv("ALLOW_ORIGINS")
	if !ok {
		return config, fmt.Errorf("ALLOW_ORIGINS not present in environment")
	}

	config.AllowOrigins = strings.Split(s, ",")

	s, ok = os.LookupEnv("SHUTDOWN_TIMEOUT")
	if !ok {
		return config, fmt.Errorf("SHUTDOWN_TIMEOUT not present in environment")
	}

	config.ShutdownTimeout, err = time.ParseDuration(s)
	if err != nil {
		return config, fmt.Errorf("fail convert SHUTDOWN_TIMEOUT: %w", err)
	}

	config.Addr, ok = os.LookupEnv("ADDR")
	if !ok {
		return config, fmt.Errorf("ADDR not present in environment")
	}

	s, ok = os.LookupEnv("GRAPHQL_QUERY_CACHE")
	if !ok {
		return config, fmt.Errorf("GRAPHQL_QUERY_CACHE not present in environment")
	}

	config.GraphQLQueryCache, err = strconv.Atoi(s)
	if err != nil {
		return config, fmt.Errorf("fail convert GRAPHQL_QUERY_CACHE: %w", err)
	}

	s, ok = os.LookupEnv("GRAPHQL_AUTOMATIC_PERSISTEND_QUERY_CACHE")
	if !ok {
		return config, fmt.Errorf("GRAPHQL_AUTOMATIC_PERSISTEND_QUERY_CACHE not present in environment")
	}

	config.GraphQLAutomaticPersistedQueryCache, err = strconv.Atoi(s)
	if err != nil {
		return config, fmt.Errorf("fail convert GRAPHQL_AUTOMATIC_PERSISTEND_QUERY_CACHE: %w", err)
	}

	s, ok = os.LookupEnv("GRAPHQL_INTROSPECTION")
	if !ok {
		return config, fmt.Errorf("GRAPHQL_INTROSPECTION not present in environment")
	}

	config.GraphQLIntrospection = s == "TRUE"

	config.GormDatabase, ok = os.LookupEnv("GORM_DATABASE")
	if !ok {
		return config, fmt.Errorf("GORM_DATABASE not present in environment")
	}

	config.GormDatabase, ok = os.LookupEnv("GORM_DSN")
	if !ok {
		return config, fmt.Errorf("GORM_DSN not present in environment")
	}

	config.SigningKey, ok = os.LookupEnv("SIGNING_KEY")
	if !ok {
		return config, fmt.Errorf("SIGNING_KEY not present in environment")
	}

	return config, nil
}

func (c *envConfiguration) state() configuration.ConfigurationState {
	builder := configuration.NewConfigurationStateBuilder()

	builder.SetAllowOrigins(c.AllowOrigins)
	builder.SetShutdownTimeout(c.ShutdownTimeout)
	builder.SetAddr(c.Addr)
	builder.SetGraphQLHandlerConfig(configuration.GraphQLHandlerConfig{
		QueryCache:                   c.GraphQLQueryCache,
		AutomaticPersistedQueryCache: c.GraphQLAutomaticPersistedQueryCache,
		Introspection:                c.GraphQLIntrospection,
	})
	builder.SetGORMRepositoryConfig(configuration.GORMRepositoryConfig{
		Database: c.GormDatabase,
		DSN:      c.GormDSN,
	})
	builder.SetSigningKey(c.SigningKey)

	return builder.Result()
}
