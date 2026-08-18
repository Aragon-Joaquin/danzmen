package flags

import (
	"flag"
	"fmt"
	"os"
)

type HelpFlag struct {
	*p_opts
}

func NewHelpFlag() FlagType {
	flag.NewFlagSet(program_help_arg, flag.ExitOnError)

	return &HelpFlag{
		p_opts: newProgramOpts(PROGRAM_HELP),
	}
}

func (*HelpFlag) UsageString() string {
	return fmt.Sprintln("Usage for ", os.Args[0], ":",
		`
 COMMANDS:
    help		This screen
    list		Output a simple screen of the tasks today
    check		(UNFINISHED) Enter in a tui to check on/off tasks
    toggle {id}		Check/uncheck a today's task
      -monthly={id}	(default)
      -long={id} 
			`)

}
