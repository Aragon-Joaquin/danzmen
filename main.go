package main

import (
	"danzmen/config"
	"danzmen/tui"
	ty "danzmen/types"
	"log"
	"time"

	"charm.land/bubbles/v2/list"

	tea "charm.land/bubbletea/v2"
)

func main() {
	cfg, err := config.ParseTOML()
	if err != nil {
		panic(err)
	}

	// op, err := ParseFlagOption()

	// NOTE: "packing" the items
	day := ty.ValidateDayOfTheWeek(time.Now().Weekday().String())
	itemsToRender := []list.Item{}

	for k, v := range cfg.Day {
		if k != day {
			continue
		}
		log.Println(k, ": ", v)
		for _, title := range v.Tasks {
			itemsToRender = append(itemsToRender, tui.CreateDZItem(title))
		}
	}

	//NOTE: start painting UI
	model := tui.CreateTUIModel(itemsToRender)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Panicf("Error running program: %e \n", err)
	}

}
