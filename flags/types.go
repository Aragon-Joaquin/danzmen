package flags

type PROGRAM_OPTION int

const (
	PROGRAM_HELP PROGRAM_OPTION = iota
	PROGRAM_LIST
	PROGRAM_CHECK
	PROGRAM_TOGGLE
)

const (
	program_help_arg   = "help"
	program_list_arg   = "list"
	program_check_arg  = "check"
	program_toggle_arg = "toggle"
)

type ProgramOpts struct {
	Type PROGRAM_OPTION
	Args []string
}
