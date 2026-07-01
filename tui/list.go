package tui

type DZList struct {
	Items      []DZItem
	selectedId int

	styles styles
}

func CreateDZList(i []DZItem, s styles) *DZList {
	//put an index for each
	for idx := range i {
		i[idx].index = idx
	}

	return &DZList{
		Items:      i,
		selectedId: 0,
		styles:     s,
	}
}

func (l *DZList) GetSelectID() int {
	return l.selectedId
}

func (l *DZList) GetSelectedItem() (item DZItem, ok bool) {
	if l.selectedId < 0 || l.selectedId >= len(l.Items) {
		return DZItem{}, false
	}

	return l.Items[l.selectedId], true
}

func (l *DZList) SetItem(idx int, item DZItem) bool {
	if idx < 0 || idx > len(l.Items) {
		return false
	}

	l.Items[idx] = item
	return true
}
