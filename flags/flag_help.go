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

type CLIFlag struct {
	Name string
	Desc string
}

// TODO: improve this
func (*HelpFlag) PrintUsage() {

	flags := []CLIFlag{
		{"help", "This screen"},
		{"list", "Output a simple screen of the tasks today"},
		{"check", "(UNFINISHED) Enter in a tui to check on/off tasks"},
		{"toggle {id}", "Check/uncheck a monthly task"},
		{"add {id} +/-{number}", "Add quantity for a monthly task "},
		{"", ""},
		{"Long modifiers: (for long tasks)", ""},
		{"toggle -l {id}", "Check/uncheck a long term task"},
		{"add -l {id} +/-{number}", "Add quantity for a longtask "},
	}

	maxLen := 0
	for _, f := range flags {
		if len(f.Name) > maxLen {
			maxLen = len(f.Name)
		}
	}

	fmt.Printf("Usage: %s\n\n", os.Args[0])
	fmt.Println("[OPTIONS]:")

	for _, f := range flags {
		fmt.Printf("  %-*s%s\n", maxLen+4, f.Name, f.Desc)
	}
	fmt.Println()
}
