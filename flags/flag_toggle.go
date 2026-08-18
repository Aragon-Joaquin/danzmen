package flags

import (
	"danzmen/db"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

type ToggleFlag struct {
	*p_opts
	Target_LongTask bool
}

var (
	flag_toggle_long_msg     = "Set to true if a Long task is getting toggle on/off, else it assumes a monthly task is going to be modified"
	flag_toggle_not_found_id = errors.New("Id not found. Does it exists in this month's tasks?")
)

func NewToggleFlag() (FlagType, error) {
	tg := flag.NewFlagSet(program_toggle_arg, flag.ExitOnError)
	if err := tg.Parse(os.Args[2:]); err != nil {
		return NewHelpFlag(), err
	}

	args := tg.Args()
	if len(args) == 0 {
		return NewHelpFlag(), errors.New("Expected an id")
	}

	target_long := tg.Bool("long", false, flag_toggle_long_msg)

	t := &ToggleFlag{
		p_opts:          newProgramOpts(PROGRAM_TOGGLE),
		Target_LongTask: *target_long,
	}

	t._args = args[:1]

	return t, nil
}

func (*ToggleFlag) UsageString() string {
	return "toggle {id}\t\tCheck/uncheck a today's task\n" +
		"  -monthly={id}\t(default)\n" +
		"  -long={id}"
}

func (p *ToggleFlag) GetID() (int, error) {
	f := p._args[0]
	id, err := strconv.Atoi(f)
	if err != nil {
		return 0, errors.New("Invalid {id}. Only accepting positive numbers.")
	}

	return id, nil
}

func (p *ToggleFlag) FlagToggle(db *db.SqliteDB, dbMonthly []*db.DBJoin_Monthly, dbLong []*db.DBLong_Tasks) error {
	if p.Target_LongTask {
		return p.flag_toggle_longTask(db, dbLong)
	}

	return p.flag_toggle_monthly(db, dbMonthly)
}

// WARN: private
func (p *ToggleFlag) flag_toggle_monthly(db *db.SqliteDB, dbTasks []*db.DBJoin_Monthly) error {
	id, err := p.GetID()
	if err != nil {
		return err
	}

	for _, v := range dbTasks {
		if v.DBMonthly_Task.Id != id {
			continue
		}

		if err := db.UpdateCompletedMonthlyTask(id, !v.Completed_At.Valid); err != nil {
			return fmt.Errorf("Error: ", err.Error())
		}

		return nil
	}
	return flag_toggle_not_found_id
}

func (p *ToggleFlag) flag_toggle_longTask(db *db.SqliteDB, dbTasks []*db.DBLong_Tasks) error {
	id, err := p.GetID()
	if err != nil {
		return err
	}

	for _, v := range dbTasks {
		if v.Id != id {
			continue
		}

		if err := db.UpdateCompletedMonthlyTask(id, !v.Completed_At.Valid); err != nil {
			return fmt.Errorf("Error: ", err.Error())
		}

		return nil
	}
	return flag_toggle_not_found_id
}
