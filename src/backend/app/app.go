package app

import (
	"iu7-2022-sd-labs/configuration"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	logger                 *log.Logger
	state                  AppState
	interrupted            bool
	interruptChan          chan os.Signal
	unsubConfigurationChan chan struct{}
}

type AppState interface {
	Update(config *configuration.ConfigurationState) (AppState, error)
	ChangeTo(state AppState) error
	Run() error
	Stop()
}

func (app *App) subscribeConfiguration(configuration *configuration.Configuration) {
	app.logger.Println("subscribe configuration")
	ch, id := configuration.Subscribe()

	go func() {
	WatchLoop:
		for {
			select {
			case config, ok := <-ch:
				if !ok {
					break WatchLoop
				}
				app.logger.Println("new configuration")

				nextState, err := app.state.Update(&config)

				if err != nil {
					app.logger.Printf("error in sub for configuration: next state update: %v\n", err)
				} else if err := app.state.ChangeTo(nextState); err != nil {
					app.logger.Printf("error in sub for configuration: state change to: %v\n", err)
				} else {
					app.logger.Print("new configuration successfuly applied\n")
				}
			case <-app.unsubConfigurationChan:
				break WatchLoop
			}
		}

		app.logger.Println("unsubscribe configuration")
		configuration.Unsubscribe(id)
	}()
}

func (app *App) subscribeInterrupt() {
	go func() {
		for range app.interruptChan {
			if app.interrupted {
				continue
			}

			app.unsubConfigurationChan <- struct{}{}
			app.state.Stop()
		}
	}()
}

func NewApp(state AppState, configuration *configuration.Configuration) (App, error) {
	app := App{
		logger:                 configuration.Logger(),
		state:                  state,
		interruptChan:          make(chan os.Signal, 1),
		unsubConfigurationChan: make(chan struct{}, 1),
	}

	signal.Notify(app.interruptChan, syscall.SIGINT, syscall.SIGTERM)
	app.subscribeConfiguration(configuration)
	app.subscribeInterrupt()

	return app, nil
}

func (app *App) Stop() {
	app.interruptChan <- syscall.SIGINT
}

func (app *App) Run() error {
	return app.state.Run()
}
