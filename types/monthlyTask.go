package types

type MonthlyTasksCfg struct {
	Name   string  `toml:"name"`
	Times  float64 `toml:"times"`
	Metric string  `toml:"metric"`
}
