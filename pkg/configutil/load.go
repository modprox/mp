package configutil

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

func GetConfigFilename(args []string) (string, error) {
	if len(args) != 2 {
		return "", errors.Errorf("expected 1 argument, got %d", len(args)-1)
	}

	return args[1], nil
}

func LoadConfig(filename string, destination interface{}) error {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "could not read config file")
	}

	if err := json.Unmarshal(bs, &destination); err != nil {
		return errors.Wrap(err, "could not parse config file")
	}

	return nil
}
