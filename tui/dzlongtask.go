package tui

import (
	"danzmen/db"
	ty "danzmen/types"
)

type DZLongTask interface {
	DZTask

	EndsOn() string
	Priority() ty.PRIORITY_TYPES
	CompletedAt() string
}

type longTask struct {
	*task
	ends         string
	priority     ty.PRIORITY_TYPES
	completed_at string
}

func CreateMultipleDZLongTask(d ...*db.DBLong_Tasks) []DZLongTask {
	var dzlong = []DZLongTask{}
	for _, v := range d {
		i := &task{
			id:        v.Id,
			title:     v.Name,
			completed: v.Completed_at.Valid,
		}

		dzlong = append(dzlong, &longTask{
			task:         i,
			ends:         v.Expires_in.String,
			priority:     v.Priority,
			completed_at: v.Completed_at.String,
		})

	}

	return dzlong
}

func (l *longTask) EndsOn() string              { return l.ends }
func (l *longTask) Priority() ty.PRIORITY_TYPES { return l.priority }
func (l *longTask) CompletedAt() string         { return l.completed_at }
