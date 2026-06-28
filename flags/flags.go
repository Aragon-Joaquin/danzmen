package flags

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func ParseOptions() (*ProgramOpts, error) {
	if len(os.Args) < 2 {
		return printHelp(), errors.New("Insufficient args provided")
	}

	switch os.Args[1] {
	case program_help_arg:
		flag.NewFlagSet(program_help_arg, flag.ExitOnError)
		return printHelp(), nil

	case program_list_arg:
		flag.NewFlagSet(program_list_arg, flag.ExitOnError)
		return &ProgramOpts{Type: PROGRAM_LIST}, nil

	case program_check_arg:
		flag.NewFlagSet(program_check_arg, flag.ExitOnError)
		return &ProgramOpts{Type: PROGRAM_CHECK}, nil

	case program_toggle_arg:
		tg := flag.NewFlagSet(program_toggle_arg, flag.ExitOnError)
		if err := tg.Parse(os.Args[2:]); err != nil {
			return printHelp(), err
		}

		args := tg.Args()
		if len(args) == 0 {
			return printHelp(), errors.New("Expected an id")
		}

		return &ProgramOpts{Type: PROGRAM_TOGGLE, Args: args[:1]}, nil

	default:
		return printHelp(), nil
	}
}

func printHelp() *ProgramOpts {
	fmt.Println(`Usage for `, os.Args[0], `: 
	- [program] help		This screen
	- [program] list		Output a simple screen of the tasks today
	- [program] check		Enter in a tui to check on/off tasks
	- [program] toggle {id}		Check/uncheck a today's task
	`)
	return &ProgramOpts{Type: PROGRAM_HELP}

}
