package tui

import (
	"danzmen/db"
	ty "danzmen/types"
	"image/color"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
)

type DZLongTask interface {
	DZTask

	ExpiresIn() string
	Priority() ty.PRIORITY_TYPES
	CompletedAt() string

	//helpers
	MM_DD_YYYY_Format() (time.Time, error)
	EndsOnXHours() float64
	RenderPriority() string
	PriorityBGColor() color.Color
}

type longTask struct {
	*task
	expires_in   string
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
			expires_in:   v.Expires_in.String,
			priority:     v.Priority,
			completed_at: v.Completed_at.String,
		})

	}

	return dzlong
}

func (l *longTask) ExpiresIn() string           { return l.expires_in }
func (l *longTask) Priority() ty.PRIORITY_TYPES { return l.priority }
func (l *longTask) CompletedAt() string         { return l.completed_at }
func (l *longTask) MM_DD_YYYY_Format() (time.Time, error) {
	return time.Parse(ty.MM_DD_YYYY, l.expires_in)
}

func (l *longTask) EndsOnXHours() float64 {
	now := time.Now()

	t, err := l.MM_DD_YYYY_Format()
	if err != nil {
		return -1
	}

	return t.Sub(now).Hours()
}

func (l *longTask) RenderPriority() string {
	return strings.ToUpper(string(l.priority))
}

func (l *longTask) PriorityBGColor() color.Color {
	switch l.priority {
	case ty.PRIO_LOW:
		return lipgloss.BrightBlack
	case ty.PRIO_MED:
		return lipgloss.Yellow
	case ty.PRIO_HIGH:
		return lipgloss.Red
	}

	return lipgloss.Black
}
