package flags

import "flag"

type ListFlag struct {
	*p_opts
}

func NewListFlag() FlagType {
	flag.NewFlagSet(program_list_arg, flag.ExitOnError)

	return &ListFlag{
		p_opts: newProgramOpts(PROGRAM_LIST),
	}
}

func (*ListFlag) UsageString() string { return "list\t\tOutput a simple screen of the tasks today" }
