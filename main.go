package main

import (
	"context"
	"danzmen/config"
	"danzmen/db"
	"danzmen/flags"
	"danzmen/tui"
	"flag"
	"fmt"
	"os"
	"time"

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

	var page int64 = 1
	if f, ok := f.(*flags.ListFlag); ok {
		page = f.ReturnPage()
	}

	// fill up the values with db
	mt, ltt, err := query_both_tasks_from_db(sdb, cfg, page)
	if err != nil {
		fmt.Printf("Couldn't query from the db %s\n", err.Error())
		os.Exit(1)
	}

	switch f.GetType() {
	case flags.PROGRAM_TOGGLE:
		if err := f.(*flags.ToggleFlag).FlagToggle(sdb, mt, ltt); err != nil {
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

	render_mt, render_ltt := format_both_tasks_for_render(mt, ltt)

	//NOTE: start painting UI
	if _, ok := f.(*flags.ListFlag); ok {
		w, h, err := xterm.GetSize(os.Stdout.Fd())
		if err != nil {
			fmt.Printf("UI Error: %s", err.Error())
			os.Exit(1)
		}

		s := tui.RenderList(render_mt, render_ltt, w, h, sdb.GetDisplayData(context.Background(), time.Now()))
		os.Stdout.WriteString(tui.CONTAINER_FOR_TOGGLE.Render(s))
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
