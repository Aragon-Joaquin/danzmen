package tui

type DZItem struct {
	title string
}

func CreateDZItem(t string) DZItem {
	return DZItem{
		title: t,
	}
}

func (i DZItem) Title() string       { return i.title }
func (i DZItem) Description() string { return i.title }
func (i DZItem) FilterValue() string { return i.title }
