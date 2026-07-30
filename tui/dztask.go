package tui

import (
	"fmt"
)

type DZTask interface {
	ID() int
	Title() string
	Completed() bool

	TitleEllipsis(length int) string
	ReturnCheckboxString() string
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
