package types

type GenericScheduleCfg struct {
	Start string   `toml:"start"`
	End   string   `toml:"end"`
	Tasks []string `toml:"tasks"`
}
