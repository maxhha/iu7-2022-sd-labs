package server

import (
	"iu7-2022-sd-labs/app"
	"iu7-2022-sd-labs/configuration"
)

func NewServerApp(config *configuration.Configuration) (app.App, error) {
	state, err := NewServerState(config)
	if err != nil {
		return app.App{}, Wrap(err, "new server state")
	}

	a, err := app.NewApp(state, config)
	if err != nil {
		return app.App{}, Wrap(err, "new app")
	}

	return a, nil
}
