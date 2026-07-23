package tui

import (
	"danzmen/db"
	ty "danzmen/types"
	"fmt"
)

type DZTask interface {
	ID() int
	Title() string
	Completed() bool

	TitleEllipsis(length int) string
	ReturnCheckboxString() string
}

func CreateMultipleDZTask(d ...*db.DBJoin_Daily) []DZTask {
	var dzitem = []DZTask{}
	for _, v := range d {
		dzitem = append(dzitem, &task{
			id:        v.DBDaily_Task.Id,
			title:     v.DBDaily_Task.Name,
			completed: ty.DBIntToBool(v.Completed),
		})
	}
	return dzitem
}

type task struct {
	id        int
	title     string
	completed bool
}

func (l *task) ID() int         { return l.id }
func (l *task) Title() string   { return l.title }
func (l *task) Completed() bool { return l.completed }

func (i *task) ReturnCheckboxString() string {
	var checked = " "
	if i.completed {
		checked = "x"
	}

	return fmt.Sprintf("[%s]", checked)
}
func (i *task) TitleEllipsis(length int) string {
	if length >= len(i.title) {
		return i.title
	}
	return fmt.Sprintf("%s...", i.title[:length])
}
