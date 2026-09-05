package main

import (
	"danzmen/config"
	"danzmen/db"
	"danzmen/tui"
)

func query_both_tasks_from_db(sdb *db.SqliteDB, cfg *config.Cfg, pageNumb int64) ([]*db.DBJoin_Monthly, []*db.DBLong_Tasks, error) {
	var err error
	mt := []*db.DBJoin_Monthly{}

	monthlyNames := cfg.GetMonthlyTasks()

	if len(monthlyNames) > 0 {
		if mt, err = sdb.CreateIfNotExistsMonthlyTasks(monthlyNames, pageNumb); err != nil {
			return nil, nil, err
		}
	}

	ltt := []*db.DBLong_Tasks{}
	ltt, err = sdb.InsertOrSelectLongTermTasks(cfg.GetLongTermTasks(), pageNumb)
	if err != nil {
		return nil, nil, err
	}

	return mt, ltt, nil
}

func format_both_tasks_for_render(mt []*db.DBJoin_Monthly, ltt []*db.DBLong_Tasks) ([]tui.DZMonthlyTask, []tui.DZLongTask) {
	monthlyToRender := []tui.DZMonthlyTask{}
	if len(mt) > 0 {
		for _, v := range tui.CreateMultipleDZMonthlyTask(mt...) {
			monthlyToRender = append(monthlyToRender, v)
		}
	}

	longToRender := []tui.DZLongTask{}
	for _, v := range tui.CreateMultipleDZLongTask(ltt...) {
		longToRender = append(longToRender, v)
	}

	return monthlyToRender, longToRender
}
