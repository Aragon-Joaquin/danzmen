package flags

import (
	"os"
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

	default:
		return NewHelpFlag(), nil
	}
}
