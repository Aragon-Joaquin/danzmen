package flags

import (
	"danzmen/db"
	"flag"
	"strconv"
)

type AddFlag struct {
	*p_opts
	Target_LongTask bool
	Id              int
	QuantityToAdd   float64
}

func NewAddFlag() (FlagType, error) {
	af := flag.NewFlagSet(program_add_arg, flag.ExitOnError)
	target_long := af.Bool("l", false, flag_modify_long_msg)

	args, err := SetArgs(af, 1)
	if err != nil {
		return NewHelpFlag(), err
	}

	t := &AddFlag{
		p_opts:          newProgramOpts(PROGRAM_ADD),
		Target_LongTask: *target_long,
	}

	id, err := GetID(args[0])
	if err != nil {
		return nil, err
	}

	q, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return nil, err
	}

	t._args = args
	t.Id = id
	t.QuantityToAdd = q

	return t, nil
}

func (p *AddFlag) FlagAddQuantity(db *db.SqliteDB) error {
	if p.Target_LongTask {
		return db.AddQuantityToLongTask(p.Id, p.QuantityToAdd)
	}
	return db.AddQuantityToMonthlyTask(p.Id, p.QuantityToAdd)
}
