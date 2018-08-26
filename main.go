package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/modprox/libmodprox/loggy"

	"github.com/pkg/errors"

	"github.com/modprox/modprox-proxy/internal/service"
)

func main() {
	log := loggy.New("modprox-proxy")
	log.Infof("--- starting up ---")

	configFilename, err := getConfigFilename(os.Args)
	if err != nil {
		log.Errorf("failed to startup: %v", err)
		os.Exit(1)
	}
	log.Infof("loading configuration from: %s", configFilename)

	config, err := loadConfig(configFilename)
	if err != nil {
		log.Errorf("failed to startup: %v", err)
		os.Exit(1)
	}
	log.Infof("starting with configuration: %s", config)

	service.NewProxy(config).Run()
}

func getConfigFilename(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.Errorf("expected 1 argument, got %d", len(args))
	}
	return args[1], nil
}

func loadConfig(filename string) (service.Configuration, error) {
	var config service.Configuration

	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, errors.Wrap(err, "could not read config file")
	}

	if err := json.Unmarshal(bs, &config); err != nil {
		return config, errors.Wrap(err, "could not parse config file")
	}

	return config, nil
}
