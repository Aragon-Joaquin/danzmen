package config

import (
	ty "danzmen/types"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Cfg struct {
	//general cfg
	Start string `toml:"start"`

	//specific cfg
	Day  map[ty.DAYS_PER_WEEK]GenericScheduleCfg `toml:"day"`
	Loop map[string]GenericScheduleCfg           `toml:"loop"`
}

const (
	USER_CONFIG_PATH = "danzmen/config.toml"
)

func (_ *Cfg) getConfigPath() string {
	home, err := os.UserConfigDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, USER_CONFIG_PATH)
}

func ParseTOML() (*Cfg, error) {
	c := generateDefaultCFGFile()
	path := c.getConfigPath()

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("Config file not specified, using default: ", err.Error())
			return c, nil
		}
		return nil, err
	}

	if _, err := toml.DecodeFile(path, c); err != nil {
		return nil, err
	}

	return c, nil
}
