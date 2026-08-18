package flags

import "flag"

type CheckFlag struct {
	*p_opts
}

func NewCheckFlag() FlagType {
	flag.NewFlagSet(program_check_arg, flag.ExitOnError)

	return &CheckFlag{
		p_opts: newProgramOpts(PROGRAM_CHECK),
	}
}

func (_ *CheckFlag) UsageString() string {
	return "check\t\t(UNFINISHED) Enter in a tui to check on/off tasks"
}
