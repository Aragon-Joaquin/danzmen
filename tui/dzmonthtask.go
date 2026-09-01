package tui

import (
	"danzmen/db"
)

type DZMonthlyTask interface {
	DZTask

	Metric() string
	TimesTotal() float64
	CurrentProgress() float64
}

type monthlyTask struct {
	*task
	metric           string
	times_of_total   float64
	current_progress float64
}

func CreateMultipleDZMonthlyTask(d ...*db.DBJoin_Monthly) []DZMonthlyTask {
	var dzitem = []DZMonthlyTask{}
	for _, v := range d {
		i := &task{
			id:        v.DBMonthly_Task.Id,
			title:     v.DBMonthly_Task.Name,
			completed: v.Completed_At.Valid}

		dzitem = append(dzitem, &monthlyTask{
			task:             i,
			metric:           v.Metric,
			times_of_total:   v.Times,
			current_progress: v.DBMonthly_Record.Times_Done.Float64,
		})
	}
	return dzitem
}

func (m *monthlyTask) Metric() string           { return m.metric }
func (m *monthlyTask) TimesTotal() float64      { return m.times_of_total }
func (m *monthlyTask) CurrentProgress() float64 { return m.current_progress }
