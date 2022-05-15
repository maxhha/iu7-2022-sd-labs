package file_configuration

import (
	"fmt"
	"iu7-2022-sd-labs/configuration"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type fileConfiguration struct {
	DebounceTime    time.Duration `yaml:"debounceTime"`
	AllowOrigins    []string      `yaml:"allowOrigins"`
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout"`
	Addr            string        `yaml:"addr"`
	Graphql         struct {
		QueryCache                   int  `yaml:"queryCacheSize"`
		AutomaticPersistedQueryCache int  `yaml:"automaticPersistedQueryCacheSize"`
		Introspection                bool `yaml:"introspection"`
	} `yaml:"graphql"`
	Gorm struct {
		Database string `yaml:"database"`
		DSN      string `yaml:"dsn"`
	}
	SigningKey string `yaml:"signingKey"`
}

func NewFileConfigurationSource(logger *log.Logger, filename string) (<-chan configuration.ConfigurationState, func(), error) {
	config, err := readConfiguration(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("read configuration: %w", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, fmt.Errorf("new watcher: %w", err)
	}

	if err = watcher.Add(filename); err != nil {
		watcher.Close()
		return nil, nil, fmt.Errorf("watcher add: %w", err)
	}

	ch := make(chan configuration.ConfigurationState, 1)
	done := make(chan struct{})
	ch <- config.state()

	go func() {
		defer close(ch)
		defer watcher.Close()

		debounceTime := config.DebounceTime
		timer := time.NewTimer(debounceTime)
		<-timer.C

		for {
			time.NewTimer(time.Millisecond)

			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					timer.Reset(debounceTime)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Printf("watcher error: %v\n", err)

			case <-timer.C:
				logger.Printf("read config file\n")
				config, err := readConfiguration(filename)
				if err != nil {
					logger.Printf("read configuration: %v\n", err)
				}

				debounceTime = config.DebounceTime
				ch <- config.state()

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

func readConfiguration(filename string) (fileConfiguration, error) {
	config := fileConfiguration{}
	file, err := os.Open(filename)
	if err != nil {
		return config, fmt.Errorf("os open: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, fmt.Errorf("yaml decoder: %w", err)
	}

	return config, nil
}

func (c *fileConfiguration) state() configuration.ConfigurationState {
	builder := configuration.NewConfigurationStateBuilder()

	builder.SetAllowOrigins(c.AllowOrigins)
	builder.SetShutdownTimeout(c.ShutdownTimeout)
	builder.SetAddr(c.Addr)
	builder.SetGraphQLHandlerConfig(configuration.GraphQLHandlerConfig{
		QueryCache:                   c.Graphql.QueryCache,
		AutomaticPersistedQueryCache: c.Graphql.AutomaticPersistedQueryCache,
		Introspection:                c.Graphql.Introspection,
	})
	builder.SetGORMRepositoryConfig(configuration.GORMRepositoryConfig{
		Database: c.Gorm.Database,
		DSN:      c.Gorm.DSN,
	})
	builder.SetSigningKey(c.SigningKey)

	return builder.Result()
}
