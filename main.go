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

	monthlyDBTasks := []*db.DBJoin_Monthly{}
	monthlyNames := cfg.GetMonthlyTasks()

	if len(monthlyNames) > 0 {
		if monthlyDBTasks, err = sdb.CreateIfNotExistsMonthlyTasks(monthlyNames); err != nil {
			log.Fatalln("CreateIfNotExists: ", err.Error())
		}
	}

	if f.Type == flags.PROGRAM_TOGGLE {
		f.FlagToggle(sdb, monthlyDBTasks)
		return
	}

	monthlyToRender := []tui.DZTask{}
	if len(monthlyDBTasks) > 0 {
		for _, v := range tui.CreateMultipleDZTask(monthlyDBTasks...) {
			monthlyToRender = append(monthlyToRender, v)
		}
	}

	ltt, err := sdb.InsertOrSelectLongTermTasks(cfg.GetLongTermTasks())
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

		s := tui.RenderList(monthlyToRender, longToRender, w, h)
		os.Stdout.WriteString(s)
		return
	}

	model := tui.CreateTUIModel(monthlyToRender, longToRender, sdb)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Panicf("Error running program: %e \n", err)
	}
}
