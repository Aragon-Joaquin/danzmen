package flags

import (
	"flag"
	"os"
	"strconv"
)

type ListFlag struct {
	*p_opts
	page int64
}

func NewListFlag() (FlagType, error) {
	f := flag.NewFlagSet(program_list_arg, flag.ExitOnError)

	flag := &ListFlag{
		p_opts: newProgramOpts(PROGRAM_LIST),
		page:   1,
	}

	if err := f.Parse(os.Args[2:]); err != nil {
		return NewHelpFlag(), err
	}

	args := f.Args()

	if len(args) >= 1 {
		q, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			q = 1
		}

		flag.page = q
	}

	return flag, nil
}

func (l *ListFlag) ReturnPage() int64 { return l.page }
