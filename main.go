package main

import (
	"danzmen/config"
	"danzmen/db"
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
	sdb, err := db.Init()
	if err != nil {
		panic(err)
	}

	// NOTE: "packing" the items
	day := ty.ValidateDayOfTheWeek(time.Now().Weekday().String())

	var names []string
	for k, v := range cfg.Day {
		if k != day {
			continue
		}
		log.Println(k, ": ", v)
		for range v.Tasks {
			names = append(names, v.Tasks...)
		}
	}

	//call the db to obtain the values
	dbTasks := []*db.DBJoin_DateRecord_Tasks{}
	if len(names) > 0 {
		dbTasks, err = sdb.CreateIfNotExistsTasks(names)
	}

	if err != nil {
		panic(err)
	}

	itemsToRender := []list.Item{}

	if len(dbTasks) > 0 {
		for _, v := range tui.CreateMultipleDZItem(dbTasks...) {
			itemsToRender = append(itemsToRender, v)
		}
	}

	//NOTE: start painting UI

	// TODO: if the flag "list" is provided, only output by using raw lipgloss with an fmt.Println to stdout
	// or i can use these? tea.WithInput(nil), tea.WithoutRenderer()

	model := tui.CreateTUIModel(itemsToRender, sdb, tui.NewSimpleStyle())
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Panicf("Error running program: %e \n", err)
	}

}
