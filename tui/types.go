package tui

import "danzmen/db"

// different modes

type RENDER_MODE int

const (
	RENDER_TUI RENDER_MODE = iota
	RENDER_ONCE
)

type DZItem struct {
	id        int
	title     string
	completed bool
}

func CreateMultipleDZItem(d ...*db.DBJoin_DateRecord_Tasks) []DZItem {
	dzitem := make([]DZItem, len(d))
	for i, v := range d {
		dzitem[i] = DZItem{
			id:        v.DBTask.Id,
			title:     v.DBTask.Name,
			completed: v.Completed.IsCompleted(),
		}

	}
	return dzitem
}
func CreateDZItem(d *db.DBJoin_DateRecord_Tasks) DZItem {
	return DZItem{
		id:        d.DBTask.Id,
		title:     d.DBTask.Name,
		completed: d.Completed.IsCompleted(),
	}
}

func (i DZItem) Title() string       { return i.title }
func (i DZItem) FilterValue() string { return i.title }
