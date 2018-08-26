package main

import (
	"os"

	"github.com/modprox/libmodprox/configutil"
	"github.com/modprox/libmodprox/loggy"
	"github.com/modprox/modprox-proxy/internal/service"
)

func main() {
	log := loggy.New("modprox-proxy")
	log.Infof("--- starting up ---")

	configFilename, err := configutil.GetConfigFilename(os.Args)
	if err != nil {
		log.Errorf("failed to startup: %v", err)
		os.Exit(1)
	}
	log.Infof("loading configuration from: %s", configFilename)

	var config service.Configuration
	if err := configutil.LoadConfig(configFilename, &config); err != nil {
		log.Errorf("failed to startup: %v", err)
		os.Exit(1)
	}
	log.Tracef("starting with configuration: %s", config)

	service.NewProxy(config).Run()
}
