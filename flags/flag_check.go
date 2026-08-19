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
