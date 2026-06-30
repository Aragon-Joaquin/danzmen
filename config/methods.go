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
