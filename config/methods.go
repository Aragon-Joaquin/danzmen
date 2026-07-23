package config

import (
	ty "danzmen/types"
	"strings"
)

func (c *Cfg) GetTasksNonRepeatableNames() []string {
	day := ty.GetTodaysDayName()
	mapNames := map[string]struct{}{}

	for k, v := range c.Day {
		if k != day {
			continue
		}

		for _, v := range v.Tasks {
			key := strings.TrimSpace(strings.ToLower(v))
			if _, ok := mapNames[key]; ok {
				continue
			}
			mapNames[key] = struct{}{}
		}
	}

	var names = []string{}
	for v := range mapNames {
		names = append(names, v)
	}

	return names
}

func (c *Cfg) GetNonRepetableLongTermTasks() []ty.LongTermTasksCfg {
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
		m = append(m, v)
	}

	return m
}
