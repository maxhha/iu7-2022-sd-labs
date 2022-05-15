package main

import (
	"iu7-2022-sd-labs/configuration"
	"iu7-2022-sd-labs/configuration/file_configuration"
	"iu7-2022-sd-labs/server"
	"log"
	"os"
)

func main() {
	configFile, exists := os.LookupEnv("CONFIG_FILE")
	if !exists {
		log.Fatalln("environment variable CONFIG_FILE is not set")
	}

	logger := log.Default()
	source, stop, err := file_configuration.NewFileConfigurationSource(logger, configFile)
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
