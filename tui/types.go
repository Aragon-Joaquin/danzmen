package tui

type DZItem struct {
	title     string
	completed bool
}

func CreateDZItem(t string, c bool) DZItem {
	return DZItem{
		title:     t,
		completed: c,
	}
}

func (i DZItem) Title() string       { return i.title }
func (i DZItem) FilterValue() string { return i.title }
