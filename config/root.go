package config

import (
	ty "danzmen/types"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Cfg struct {
	//specific cfg
	Month map[ty.ALL_MONTHS]struct {
		Tasks []toml.Primitive `toml:"tasks"`
	} `toml:"month"`

	LongTerm struct {
		Tasks []ty.LongTermTasksCfg
	} `toml:"longterm"`
	monthParsed    []ty.MonthlyTasksCfg
	longTermParsed []ty.LongTermTasksCfg
}

func generateDefaultCFGFile() *Cfg {
	return &Cfg{
		monthParsed:    []ty.MonthlyTasksCfg{},
		longTermParsed: []ty.LongTermTasksCfg{},
	}
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
			fmt.Println("Config file not specified, using default.", err.Error())
			return c, nil
		}
		return nil, err
	}

	meta, err := toml.DecodeFile(path, c)
	if err != nil {
		terr, ok := err.(toml.ParseError)
		if !ok {
			return nil, err
		}

		return nil, errors.New(terr.ErrorWithUsage())
	}

	c.monthParsed = c.getNonRepeatableMonthlyTasks(&meta)
	c.longTermParsed = c.getNonRepetableLongTermTasks()

	for i, v := range meta.Undecoded() {
		if i == 0 {
			ty.TermSetColor(ty.Red)
		}

		fmt.Printf("- UNKNOWN FIELD: %s\n", v.String())
	}
	ty.TermResetColor()

	return c, nil
}
