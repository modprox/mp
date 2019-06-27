package main

import (
	"os"

	"oss.indeed.com/go/modprox/pkg/configutil"
	"oss.indeed.com/go/modprox/pkg/loggy"
	"oss.indeed.com/go/modprox/registry"
	"oss.indeed.com/go/modprox/registry/config"
)

// generate webpage statics
//go:generate go run github.com/shoenig/petrify/v4/cmd/petrify -prefix ../../registry -o ../../registry/static/generated.go -pkg static ../../registry/static/...

func main() {
	log := loggy.New("modprox-registry")
	log.Infof("--- starting up ---")

	configFilename, err := configutil.GetConfigFilename(os.Args)
	if err != nil {
		log.Errorf("failed to startup: %v", err)
		os.Exit(1)
	}
	log.Infof("loading configuration from: %s", configFilename)

	var configuration config.Configuration
	if err := configutil.LoadConfig(configFilename, &configuration); err != nil {
		log.Errorf("failed to startup: %v", err)
		os.Exit(1)
	}
	log.Tracef("starting with configuration: %s", configuration)

	registry.Start(configuration)
}
