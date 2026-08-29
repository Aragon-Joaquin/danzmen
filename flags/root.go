package flags

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseOptions() (FlagType, error) {
	if len(os.Args) < 2 {
		return NewHelpFlag(), nil
	}

	switch os.Args[1] {
	case program_help_arg:
		return NewHelpFlag(), nil

	case program_list_arg:
		return NewListFlag(), nil

	case program_check_arg:
		return NewCheckFlag(), nil

	case program_toggle_arg:
		return NewToggleFlag()

	case program_add_arg:
		return NewAddFlag()

	default:
		return NewHelpFlag(), nil
	}
}

func GetID(arg string) (int, error) {
	id, err := strconv.Atoi(arg)
	if err != nil {
		return 0, errors.New("Invalid {id}. Only accepting positive numbers.")
	}

	return id, nil
}

func SetArgs(f *flag.FlagSet, arguments_needed int) ([]string, error) {
	if len(os.Args) < 2 {
		return nil, errors.New("missing subcommand")
	}

	if err := f.Parse(os.Args[2:]); err != nil {
		return nil, err
	}

	args := f.Args()

	// arg 0 = program name
	// arg 1 = action (toggle/check/list...)
	// arg 2 = id

	len_arg := len(args)
	if len_arg != arguments_needed {
		builder := strings.Builder{}
		builder.WriteString("danzmen cmd [")

		for idx := range max(len_arg, arguments_needed) {
			if len_arg < arguments_needed && len_arg <= idx {
				builder.WriteString(" ... ")
				continue
			}

			if idx == arguments_needed {
				builder.WriteString("]\n\t[ EXTRA ARGS -------> ")
			}

			arg := args[idx]
			builder.WriteString(" ")
			builder.WriteString(arg)
			builder.WriteString(" ")
			continue
		}
		builder.WriteString(" ]")
		return nil, fmt.Errorf("Needed AT LEAST %d arguments and received %d:\t\n %s", arguments_needed, len(args), builder.String())

	}

	return args, nil
}
