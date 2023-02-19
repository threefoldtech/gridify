package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const CONFIG_FIlE = ".gridifyconfig"

type Config struct {
	Mnemonics string `json:"mnemonics"`
	Network   string `json:"network"`
}

func SaveConfigData(mnemonics, network string) error {

	configDir, err := os.UserConfigDir()
	if err != nil {
		return errors.Wrap(err, "could not get user configuration directory")
	}
	path := filepath.Join(configDir, CONFIG_FIlE)
	configFile, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "could not create configuration file %s", path)
	}
	defer configFile.Close()
	config := Config{
		Mnemonics: mnemonics,
		Network:   network,
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		return errors.Wrapf(err, "could not marshal configuration data %+v", config)
	}
	_, err = configFile.Write(configJSON)
	if err != nil {
		return errors.Wrapf(err, "could not wirte configuration data to file %s", configFile.Name())
	}
	return nil
}

func LoadConfigData() (Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, errors.Wrap(err, "could not get user configuration directory")
	}
	path := filepath.Join(configDir, CONFIG_FIlE)
	configJSON, err := os.ReadFile(path)
	if err != nil {
		return Config{}, errors.Wrapf(err, "could not read configuration file %s", path)
	}
	config := Config{}
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		return config, errors.Wrapf(err, "could not unmarshal configuration data %s", configJSON)
	}
	return config, nil
}
