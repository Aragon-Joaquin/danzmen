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

	//TODO: make this into a db package so then i can use the types in tui/
	_, err = initDB()
	if err != nil {
		panic(err)
	}

	// NOTE: "packing" the items
	day := ty.ValidateDayOfTheWeek(time.Now().Weekday().String())
	itemsToRender := []list.Item{
		tui.CreateDZItem("Task A", true),
		tui.CreateDZItem("Task B", false),
		tui.CreateDZItem("Task C", false),
		tui.CreateDZItem("Task ...", true),
	}

	for k, v := range cfg.Day {
		if k != day {
			continue
		}
		log.Println(k, ": ", v)
		for range v.Tasks {
			itemsToRender = append(itemsToRender, itemsToRender...)
		}
	}

	//NOTE: start painting UI
	model := tui.CreateTUIModel(itemsToRender, tui.NewSimpleStyle())
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Panicf("Error running program: %e \n", err)
	}

}
