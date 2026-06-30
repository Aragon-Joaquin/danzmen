package main

import (
	"danzmen/config"
	"danzmen/db"
	"danzmen/flags"
	"danzmen/tui"
	"log"

	"charm.land/bubbles/v2/list"

	tea "charm.land/bubbletea/v2"
)

func main() {
	//NOTE: flag parsing
	f, err := flags.ParseOptions()
	if err != nil {
		log.Fatalln(err.Error())
	}

	if f.Type == flags.PROGRAM_HELP {
		return
	}

	//NOTE: toml
	cfg, err := config.ParseTOML()
	if err != nil {
		log.Fatalln(err.Error())
	}

	names := cfg.GetTasksNonRepeatableNames()

	//NOTE: call the db to obtain the values
	sdb, err := db.Init()
	if err != nil {
		log.Fatalln(err.Error())
	}

	dbTasks := []*db.DBJoin_DateRecord_Tasks{}
	if len(names) > 0 {
		if dbTasks, err = sdb.CreateIfNotExistsTasks(names); err != nil {
			log.Fatalln(err.Error())
		}
	}

	if f.Type == flags.PROGRAM_TOGGLE {
		f.FlagToggle(sdb, dbTasks)
		return
	}

	itemsToRender := []list.Item{}
	if len(dbTasks) > 0 {
		for _, v := range tui.CreateMultipleDZItem(dbTasks...) {
			itemsToRender = append(itemsToRender, v)
		}
	}

	//NOTE: start painting UI
	var model tui.TuiModel
	switch f.Type {
	case flags.PROGRAM_CHECK:
		model = tui.CreateTUIModel(itemsToRender, sdb, tui.NewSimpleStyle(), false)
	case flags.PROGRAM_LIST:
		model = tui.CreateTUIModel(itemsToRender, sdb, tui.NewSimpleStyle(), true)
	default:
		panic("invalid option")
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Panicf("Error running program: %e \n", err)
	}
}
