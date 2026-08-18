package flags

const (
	program_help_arg   = "help"
	program_list_arg   = "list"
	program_check_arg  = "check"
	program_toggle_arg = "toggle"
)

type PROGRAM_OPTION int

const (
	PROGRAM_HELP PROGRAM_OPTION = iota
	PROGRAM_LIST
	PROGRAM_CHECK
	PROGRAM_TOGGLE
)

type FlagType interface {
	GetType() PROGRAM_OPTION
	GetArgs() []string

	UsageString() string
}

type p_opts struct {
	_type PROGRAM_OPTION
	_args []string
}

func newProgramOpts(t PROGRAM_OPTION) *p_opts {
	return &p_opts{
		_type: t,
		_args: []string{},
	}
}

func (p *p_opts) GetType() PROGRAM_OPTION { return p._type }
func (p *p_opts) GetArgs() []string       { return p._args }
