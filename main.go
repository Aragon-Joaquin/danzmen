package main

import (
	"danzmen/config"
	"danzmen/db"
	"danzmen/flags"
	"danzmen/tui"
	"log"

	tea "charm.land/bubbletea/v2"
)

func main() {
	//NOTE: flag parsing
	f, err := flags.ParseOptions()
	if err != nil {
		log.Fatalln("FLAG: ", err.Error())
	}

	if f.Type == flags.PROGRAM_HELP {
		return
	}

	//NOTE: toml
	cfg, err := config.ParseTOML()
	if err != nil {
		log.Fatalln("TOML: ", err.Error())
		return
	}

	//NOTE: call the db to obtain the values
	sdb, err := db.Init()
	if err != nil {
		log.Fatalln("DB: ", err.Error())
	}

	dbTasks := []*db.DBJoin_Daily{}
	dailyNames := cfg.GetTasksNonRepeatableNames()

	if len(dailyNames) > 0 {
		if dbTasks, err = sdb.CreateIfNotExistsTasks(dailyNames); err != nil {
			log.Fatalln("CreateIfNotExists: ", err.Error())
		}
	}

	if f.Type == flags.PROGRAM_TOGGLE {
		f.FlagToggle(sdb, dbTasks)
		return
	}

	dailyToRender := []tui.DZTask{}
	if len(dbTasks) > 0 {
		for _, v := range tui.CreateMultipleDZTask(dbTasks...) {
			dailyToRender = append(dailyToRender, v)
		}
	}

	longTerm := cfg.GetNonRepetableLongTermTasks()
	longToRender := []tui.DZLongTask{}
	if len(longTerm) > 0 {
		for _, v := range tui.CreateMultipleDZTask(dbTasks...) {
			dailyToRender = append(dailyToRender, v)
		}
	}

	//NOTE: start painting UI
	var model tui.TuiModel
	switch f.Type {
	case flags.PROGRAM_CHECK:
		model = tui.CreateTUIModel(dailyToRender, longToRender, sdb, false)
	case flags.PROGRAM_LIST:
		model = tui.CreateTUIModel(dailyToRender, longToRender, sdb, true)
	default:
		panic("invalid option")
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Panicf("Error running program: %e \n", err)
	}
}
