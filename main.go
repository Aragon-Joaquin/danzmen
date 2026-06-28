package main

import (
	"danzmen/config"
	"danzmen/db"
	"danzmen/flags"
	"danzmen/tui"
	ty "danzmen/types"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"charm.land/bubbles/v2/list"

	tea "charm.land/bubbletea/v2"
)

func main() {
	//NOTE: flag parsing
	f, err := flags.ParseOptions()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	if f.Type == flags.PROGRAM_HELP {
		return
	}

	//toml
	cfg, err := config.ParseTOML()
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
		for range v.Tasks {
			names = append(names, v.Tasks...)
		}
	}

	//NOTE: call the db to obtain the values
	sdb, err := db.Init()
	if err != nil {
		panic(err)
	}

	dbTasks := []*db.DBJoin_DateRecord_Tasks{}
	if len(names) > 0 {
		dbTasks, err = sdb.CreateIfNotExistsTasks(names)
	}

	if f.Type == flags.PROGRAM_TOGGLE {
		f := f.Args[0]
		id, err := strconv.Atoi(f)
		if err != nil {
			fmt.Println("Invalid {id}. Only accepting positive numbers.")
			os.Exit(1)
		}

		for _, v := range dbTasks {
			if v.Id == id {
				var c int = 0
				if v.Completed == 0 {
					c = 1
				}
				if err := sdb.UpdateCompletedTask(id, c); err != nil {
					fmt.Println("Error: ", err.Error())
					os.Exit(1)
					return
				}
				fmt.Println("Success")
				return
			}
		}
		fmt.Println("Id not found. Does it exists in today's date?")
		os.Exit(1)
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

	var model tui.TuiModel
	switch f.Type {
	case flags.PROGRAM_CHECK:
		model = tui.CreateTUIModel(itemsToRender, sdb, tui.NewSimpleStyle(), false)
	case flags.PROGRAM_LIST:
		model = tui.CreateTUIModel(itemsToRender, sdb, tui.NewSimpleStyle(), true)

	case flags.PROGRAM_TOGGLE:
		return

	default:
		panic("invalid option")
	}

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Panicf("Error running program: %e \n", err)
	}

}
