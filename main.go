package main

import (
	"danzmen/config"
	"danzmen/db"
	"danzmen/flags"
	"danzmen/tui"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
	xterm "github.com/charmbracelet/x/term"
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

	ltt, err := sdb.InsertOrSelectLongTermTasks(cfg.GetNonRepetableLongTermTasks())
	if err != nil {
		log.Fatalln("InsertOrSelectLongTermTasks: ", err)
	}

	longToRender := []tui.DZLongTask{}
	for _, v := range tui.CreateMultipleDZLongTask(ltt...) {
		longToRender = append(longToRender, v)
	}

	//NOTE: start painting UI
	if f.Type == flags.PROGRAM_LIST {
		w, h, err := xterm.GetSize(os.Stdout.Fd())
		if err != nil {
			log.Fatalln(err)
		}
		// sure, my terminal has 4px of padding
		// this is the worst. remake this in lua. please
		s := tui.RenderList(dailyToRender, longToRender, w-4, h)
		os.Stdout.WriteString(s)
		return
	}

	model := tui.CreateTUIModel(dailyToRender, longToRender, sdb)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Panicf("Error running program: %e \n", err)
	}
}
