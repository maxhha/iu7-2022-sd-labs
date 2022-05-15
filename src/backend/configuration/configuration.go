package configuration

import (
	"fmt"
	"log"
	"sync"
)

type Configuration struct {
	logger      *log.Logger
	stateSource <-chan ConfigurationState
	mu          sync.Mutex
	current     ConfigurationState
	subsCounter int
	subs        map[int]chan<- ConfigurationState
}

func NewConfiguration(logger *log.Logger, stateSource <-chan ConfigurationState) *Configuration {
	subs := make(map[int]chan<- ConfigurationState)

	config := Configuration{
		logger:      logger,
		stateSource: stateSource,
		subs:        subs,
		current:     <-stateSource,
	}
	go config.handleSource()

	return &config
}

func (c *Configuration) handleSource() {
	for state := range c.stateSource {
		c.mu.Lock()

		fmt.Printf("here %#v\n", state)
		c.current = state

		for _, ch := range c.subs {
			go func(ch chan<- ConfigurationState) {
				ch <- state
			}(ch)
		}

		c.mu.Unlock()
	}
}

func (c *Configuration) Current() ConfigurationState {
	return c.current
}

func (c *Configuration) Logger() *log.Logger {
	return c.logger
}

func (c *Configuration) Subscribe() (<-chan ConfigurationState, int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch := make(chan ConfigurationState, 1)
	id := c.subsCounter
	c.subsCounter++
	c.subs[id] = ch

	return ch, id
}

func (c *Configuration) Unsubscribe(id int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.subs, id)
}
