package flags

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	flag_toggle_long_msg = "Set to true if a Long task is getting toggle on/off, else it assumes a monthly task is going to be modified"
)

func ParseOptions() (*ProgramOpts, error) {
	if len(os.Args) < 2 {
		return printHelp(), nil
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

		tg.Bool("long", false, flag_toggle_long_msg)

		return &ProgramOpts{Type: PROGRAM_TOGGLE, Args: args[:1]}, nil

	default:
		return printHelp(), nil
	}
}

func printHelp() *ProgramOpts {
	fmt.Println("Usage for ", os.Args[0], ":",
		`
 COMMANDS:
    help		This screen
    list		Output a simple screen of the tasks today
    check		(UNFINISHED) Enter in a tui to check on/off tasks
    toggle {id}		Check/uncheck a today's task
      -monthly={id}	(default)
      -long={id} 
			`)
	return &ProgramOpts{Type: PROGRAM_HELP}

}
