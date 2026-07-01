package tui

import (
	"danzmen/db"
	ty "danzmen/types"
	"fmt"
)

// TODO: make DZItem an interface instead of a struct
type DZItem struct {
	id        int
	index     int
	title     string
	completed bool
}

func CreateMultipleDZItem(d ...*db.DBJoin_DateRecord_Tasks) []DZItem {
	var dzitem = []DZItem{}
	for _, v := range d {
		dzitem = append(dzitem, DZItem{
			id:        v.DBTask.Id,
			title:     v.DBTask.Name,
			completed: ty.DBIntToBool(v.Completed),
		})

	}
	return dzitem
}
func CreateDZItem(d *db.DBJoin_DateRecord_Tasks) DZItem {
	return DZItem{
		id:        d.DBTask.Id,
		title:     d.DBTask.Name,
		completed: ty.DBIntToBool(d.Completed),
	}
}

func (i DZItem) Title() string {
	return i.title
}

func (i DZItem) ReturnCheckboxString() string {
	var checked = " "
	if i.completed {
		checked = "x"
	}

	return fmt.Sprintf("[%s]", checked)
}
func (i DZItem) FilterValue() string { return i.title }
func (i DZItem) Description() string { return "" }
