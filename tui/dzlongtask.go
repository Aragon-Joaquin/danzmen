package tui

import (
	"danzmen/db"
	ty "danzmen/types"
	"fmt"
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
	HumanReadableEndsIn() string
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
			completed: v.Completed_At.Valid,
		}

		dzlong = append(dzlong, &longTask{
			task:         i,
			expires_in:   v.Expires_in.String,
			priority:     v.Priority,
			completed_at: v.Completed_At.String,
		})

	}

	return dzlong
}

func (l *longTask) ExpiresIn() string           { return l.expires_in }
func (l *longTask) Priority() ty.PRIORITY_TYPES { return l.priority }
func (l *longTask) CompletedAt() string         { return l.completed_at }
func (l *longTask) MM_DD_YYYY_Format() (time.Time, error) {
	return time.Parse(string(ty.MM_DD_YYYY), l.expires_in)
}

func (l *longTask) RenderPriority() string {
	return strings.ToUpper(string(l.priority))
}

func (l *longTask) PriorityBGColor() color.Color {
	switch l.priority {
	case ty.PRIO_LOW:
		return lipgloss.Cyan
	case ty.PRIO_MED:
		return lipgloss.Yellow
	case ty.PRIO_HIGH:
		return lipgloss.Red
	}

	return lipgloss.Black
}

const (
	hours_to_be_considered_not_that_urgent = 72
)

// more easy to follow up
func (l *longTask) HumanReadableEndsIn() string {
	h := l.endsOnXHours()

	if h == -1 {
		return "???d"
	}

	if h > hours_to_be_considered_not_that_urgent {
		return fmt.Sprintf("%.0fd", h/24)
	}

	return fmt.Sprintf("%.2fh", h)
}

// WARN: private
func (l *longTask) endsOnXHours() float64 {
	now := time.Now()

	t, err := l.MM_DD_YYYY_Format()
	if err != nil {
		return -1
	}

	return t.Sub(now).Hours()
}
