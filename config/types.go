package config

type GenericScheduleCfg struct {
	Start string   `toml:"start"`
	End   string   `toml:"end"`
	Tasks []string `toml:"tasks"`
}

func generateDefaultCFGFile() *Cfg {
	return &Cfg{
		Start: "sunday",
	}
}
