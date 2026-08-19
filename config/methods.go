package config

import (
	ty "danzmen/types"
	"log"
	"strings"

	"github.com/BurntSushi/toml"
)

func (c *Cfg) getNonRepeatableMonthlyTasks(m *toml.MetaData) []ty.MonthlyTasksCfg {
	cMonth := ty.GetTodaysMonth()
	mapNames := map[string]ty.MonthlyTasksCfg{}

	for k, v := range c.Month {
		if k != cMonth && k != ty.EVERY {
			continue
		}

		for _, v := range v.Tasks {
			var t ty.MonthlyTasksCfg
			var name string
			if err := m.PrimitiveDecode(v, &name); err == nil {
				t = ty.MonthlyTasksCfg{
					Name:   name,
					Times:  0,
					Metric: "",
				}
			} else if err := m.PrimitiveDecode(v, &t); err != nil {
				log.Fatalln("Unrecognizable monthly task value: ", v)
			}

			if t.Name == "" {
				continue
			}

			key := strings.TrimSpace(strings.ToLower(t.Name))
			if _, ok := mapNames[key]; ok {
				continue
			}

			mapNames[key] = t
		}
	}

	var names = []ty.MonthlyTasksCfg{}
	for _, v := range mapNames {
		names = append(names, v)
	}

	return names
}

func (c *Cfg) getNonRepetableLongTermTasks() []ty.LongTermTasksCfg {
	mapNames := map[string]ty.LongTermTasksCfg{}

	for _, t := range c.LongTerm.Tasks {
		if t.Name == "" {
			continue
		}

		if _, ok := mapNames[t.Name]; ok {
			continue
		}

		mapNames[t.Name] = t
	}

	m := []ty.LongTermTasksCfg{}
	for _, v := range mapNames {
		//validate if they're correctly parsed
		if err := v.ValidateExpires_In(); err != nil {
			log.Fatalln(err)
			return nil
		}

		if err := v.ValidatePriority(); err != nil {
			log.Fatalln(err)
			return nil
		}

		m = append(m, v)
	}

	return m
}

func (c *Cfg) GetMonthlyTasks() []ty.MonthlyTasksCfg   { return c.monthParsed }
func (c *Cfg) GetLongTermTasks() []ty.LongTermTasksCfg { return c.longTermParsed }
