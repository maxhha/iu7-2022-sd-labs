package main

import (
	"iu7-2022-sd-labs/configuration"
	"iu7-2022-sd-labs/configuration/env_configuration"
	"iu7-2022-sd-labs/configuration/file_configuration"
	"iu7-2022-sd-labs/server"
	"log"
	"os"
)

func main() {
	var (
		source <-chan configuration.ConfigurationState
		stop   func()
		err    error
	)
	logger := log.Default()
	configFile, exists := os.LookupEnv("CONFIG_FILE")

	if exists {
		logger.Println("CONFIG_FILE present in env. Use file configuration source.")
		source, stop, err = file_configuration.NewFileConfigurationSource(logger, configFile)
	} else {
		source, stop, err = env_configuration.NewEnvConfigurationSource(logger)
	}

	if err != nil {
		panic(err)
	}

	defer stop()

	config := configuration.NewConfiguration(logger, source)
	app, err := server.NewServerApp(config)
	if err != nil {
		panic(err)
	}

	app.Run()
}
