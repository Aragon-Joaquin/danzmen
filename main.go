package main

import (
	"danzmen/config"
	"danzmen/db"
	"danzmen/flags"
	"danzmen/tui"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	xterm "github.com/charmbracelet/x/term"
)

func main() {
	//NOTE: flag parsing
	f, err := flags.ParseOptions()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if f.Type == flags.PROGRAM_HELP {
		return
	}

	//NOTE: toml
	cfg, err := config.ParseTOML()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//NOTE: call the db to obtain the values
	sdb, err := db.Init()
	if err != nil {
		fmt.Printf("Error while saving to the Database: %s", err.Error())
		os.Exit(1)
	}

	monthlyDBTasks := []*db.DBJoin_Monthly{}
	monthlyNames := cfg.GetMonthlyTasks()

	if len(monthlyNames) > 0 {
		if monthlyDBTasks, err = sdb.CreateIfNotExistsMonthlyTasks(monthlyNames); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	if f.Type == flags.PROGRAM_TOGGLE {
		f.FlagToggle(sdb, monthlyDBTasks)
		return
	}

	monthlyToRender := []tui.DZMonthlyTask{}
	if len(monthlyDBTasks) > 0 {
		for _, v := range tui.CreateMultipleDZMonthlyTask(monthlyDBTasks...) {
			monthlyToRender = append(monthlyToRender, v)
		}
	}

	ltt, err := sdb.InsertOrSelectLongTermTasks(cfg.GetLongTermTasks())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	longToRender := []tui.DZLongTask{}
	for _, v := range tui.CreateMultipleDZLongTask(ltt...) {
		longToRender = append(longToRender, v)
	}

	//NOTE: start painting UI
	if f.Type == flags.PROGRAM_LIST {
		w, h, err := xterm.GetSize(os.Stdout.Fd())
		if err != nil {
			fmt.Printf("UI Error: %s", err.Error())
			os.Exit(1)
		}

		s := tui.RenderList(monthlyToRender, longToRender, w, h)
		c := tui.CONTAINER_FOR_TOGGLE
		os.Stdout.WriteString(c.Render(s))
		return
	}

	model := tui.CreateTUIModel(monthlyToRender, longToRender, sdb)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %s", err.Error())
		os.Exit(1)
	}
}
