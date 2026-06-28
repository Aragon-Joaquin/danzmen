package tui

import (
	"danzmen/db"
	"fmt"
)

func DBIntToBool(i int) bool {
	if i == 1 {
		return true
	}
	return false
}

type DZItem struct {
	id        int
	title     string
	completed bool
}

func CreateMultipleDZItem(d ...*db.DBJoin_DateRecord_Tasks) []DZItem {
	var dzitem = []DZItem{}
	for _, v := range d {
		dzitem = append(dzitem, DZItem{
			id:        v.DBTask.Id,
			title:     v.DBTask.Name,
			completed: DBIntToBool(v.Completed),
		})

	}
	return dzitem
}
func CreateDZItem(d *db.DBJoin_DateRecord_Tasks) DZItem {
	return DZItem{
		id:        d.DBTask.Id,
		title:     d.DBTask.Name,
		completed: DBIntToBool(d.Completed),
	}
}

func (i DZItem) Title() string {
	var checked = " "
	if i.completed {
		checked = "x"
	}

	return fmt.Sprintf("[%s] %s", checked, i.title)
}
func (i DZItem) FilterValue() string { return i.title }
func (i DZItem) Description() string { return "" }
