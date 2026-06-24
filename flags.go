package main

import (
	"flag"
	"log"
	"os"
)

type PROGRAM_OPTION int

const (
	PROGRAM_HELP PROGRAM_OPTION = iota
	PROGRAM_LIST
	PROGRAM_CHECK
)

const (
	program_help_arg  = "help"
	program_list_arg  = "list"
	program_check_arg = "check"
)

func ParseFlagOption() (PROGRAM_OPTION, error) {
	if len(os.Args) < 2 {
		log.Println("Insufficient args provided")
		printHelp()
		os.Exit(1)
	}

	//or flag.ContinueOnError idk.
	// dont make this an array, maybe i would need the return types
	flag.NewFlagSet(program_help_arg, flag.ExitOnError)
	flag.NewFlagSet(program_check_arg, flag.ExitOnError)
	flag.NewFlagSet(program_list_arg, flag.ExitOnError)

	switch os.Args[1] {
	case program_help_arg:
		return printHelp(), nil

	case program_list_arg:
		return PROGRAM_LIST, nil

	case program_check_arg:
		return PROGRAM_CHECK, nil

	default:
		return printHelp(), nil
	}
}

func printHelp() PROGRAM_OPTION {
	log.Println(`
	Usage for `, os.Args[0], `: 
	- [program] help:				This screen
	- [program] list:				Output a simple screen of the tasks today
	- [program] check:			Enter in a tui to check on/off tasks
	`)
	return PROGRAM_HELP

}
