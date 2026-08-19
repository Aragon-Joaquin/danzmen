package main

import (
	"danzmen/config"
	"danzmen/db"
	"danzmen/flags"
	"danzmen/tui"
	"flag"
	"fmt"
	"os"

	xterm "github.com/charmbracelet/x/term"
)

func main() {
	//NOTE: flag parsing
	flag.Parse()
	f, err := flags.ParseOptions()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if f.GetType() == flags.PROGRAM_HELP {
		f.(*flags.HelpFlag).PrintUsage()
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

	// fill up the values with db
	monthlyDBTasks := []*db.DBJoin_Monthly{}
	monthlyNames := cfg.GetMonthlyTasks()

	if len(monthlyNames) > 0 {
		if monthlyDBTasks, err = sdb.CreateIfNotExistsMonthlyTasks(monthlyNames); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	ltt, err := sdb.InsertOrSelectLongTermTasks(cfg.GetLongTermTasks())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	switch f.GetType() {
	case flags.PROGRAM_TOGGLE:
		if err := f.(*flags.ToggleFlag).FlagToggle(sdb, monthlyDBTasks, ltt); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	case flags.PROGRAM_ADD:
		if err := f.(*flags.AddFlag).FlagAddQuantity(sdb); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	}

	//create render-able objects
	monthlyToRender := []tui.DZMonthlyTask{}
	if len(monthlyDBTasks) > 0 {
		for _, v := range tui.CreateMultipleDZMonthlyTask(monthlyDBTasks...) {
			monthlyToRender = append(monthlyToRender, v)
		}
	}

	longToRender := []tui.DZLongTask{}
	for _, v := range tui.CreateMultipleDZLongTask(ltt...) {
		longToRender = append(longToRender, v)
	}

	//NOTE: start painting UI
	if _, ok := f.(*flags.ListFlag); ok {
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

	//else, it checked the "CHECK" flag
	// TODO: make this work lol
	fmt.Println("	>> Check mode isn't complete yet")

	// model := tui.CreateTUIModel(monthlyToRender, longToRender, sdb)
	// p := tea.NewProgram(model)
	// if _, err := p.Run(); err != nil {
	// 	fmt.Printf("Error running program: %s", err.Error())
	// 	os.Exit(1)
	// }
}
