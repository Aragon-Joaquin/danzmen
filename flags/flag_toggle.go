package flags

import (
	"danzmen/db"
	"errors"
	"flag"
	"fmt"
)

type ToggleFlag struct {
	*p_opts
	Target_LongTask bool
	Id              int
}

var (
	flag_modify_long_msg     = "Set to true if a Long task is going to be modified, else it assumes it's a monthly task."
	flag_toggle_not_found_id = errors.New("Id not found. Does it exists in this month's tasks?")
)

func NewToggleFlag() (FlagType, error) {
	tg := flag.NewFlagSet(program_toggle_arg, flag.ExitOnError)
	target_long := tg.Bool("l", false, flag_modify_long_msg)

	args, err := SetArgs(tg, 1)
	if err != nil {
		return NewHelpFlag(), err
	}

	t := &ToggleFlag{
		p_opts:          newProgramOpts(PROGRAM_TOGGLE),
		Target_LongTask: *target_long,
	}

	id, err := GetID(args[0])
	if err != nil {
		return nil, err
	}

	t._args = args
	t.Id = id

	return t, nil
}

func (p *ToggleFlag) FlagToggle(db *db.SqliteDB, dbMonthly []*db.DBJoin_Monthly, dbLong []*db.DBLong_Tasks) error {
	if p.Target_LongTask {
		return p.flag_toggle_longTask(db, dbLong)
	}

	return p.flag_toggle_monthly(db, dbMonthly)
}

// WARN: private
func (p *ToggleFlag) flag_toggle_monthly(db *db.SqliteDB, dbTasks []*db.DBJoin_Monthly) error {
	for _, v := range dbTasks {
		if v.DBMonthly_Task.Id != p.Id {
			continue
		}

		if err := db.UpdateCompletedMonthlyTask(p.Id, !v.Completed_At.Valid); err != nil {
			return fmt.Errorf("Error: %s", err.Error())
		}

		return nil
	}
	return flag_toggle_not_found_id
}

func (p *ToggleFlag) flag_toggle_longTask(db *db.SqliteDB, dbTasks []*db.DBLong_Tasks) error {
	for _, v := range dbTasks {
		if v.Id != p.Id {
			continue
		}

		if err := db.UpdateCompletedLongTask(p.Id, !v.Completed_At.Valid); err != nil {
			return fmt.Errorf("Error: %s", err.Error())
		}

		return nil
	}
	return flag_toggle_not_found_id
}
