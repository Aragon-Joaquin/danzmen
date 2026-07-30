package tui

import (
	"cmp"
	"fmt"
	"slices"

	"charm.land/lipgloss/v2"
)

type DZList interface {
	GetSelectID() int
	SelectedItem() (item DZTask, ok bool)
	SetItem(idx int, item DZTask) bool
	Counts() (total int, completed int)
	CastToLongTask(i DZTask) (t DZLongTask, ok bool)

	//ui things
	SetHeight(h int)
	SetWidth(w int)
	SetSizes(w, h int)
	View() string
}

type listModel struct {
	h, w       int
	items      []listItem
	selectedId int

	styles styles
}

type listItem struct {
	item DZTask
	id   int
}

func CreateDZLongList(items []DZLongTask, s styles, w, h int) DZList {
	var listTasks = []listItem{}
	for idx, v := range items {
		listTasks = append(listTasks, listItem{
			item: v,
			id:   idx,
		})
	}

	return &listModel{
		items:      listTasks,
		selectedId: 0,
		styles:     s,
		h:          h,
		w:          w,
	}
}

func CreateDZList(i []DZMonthlyTask, s styles, w, h int) DZList {
	//put an index for each
	var listTasks = []listItem{}
	for idx, v := range i {
		listTasks = append(listTasks, listItem{
			item: v,
			id:   idx,
		})
	}

	return &listModel{
		items:      listTasks,
		selectedId: 0,
		styles:     s,
		h:          h,
		w:          w,
	}
}

func (l *listModel) CastToLongTask(i DZTask) (t DZLongTask, ok bool) {
	t, ok = i.(DZLongTask)
	return
}

func (l *listModel) GetSelectID() int {
	return l.selectedId
}

func (l *listModel) SelectedItem() (item DZTask, ok bool) {
	if l.selectedId < 1 || l.selectedId >= len(l.items) {
		return nil, false
	}

	i := l.items[l.selectedId]
	return i.item, true
}

func (l *listModel) SetItem(idx int, item DZTask) bool {
	if idx < 0 || idx > len(l.items) {
		return false
	}

	lTask := &listItem{
		item: item,
		id:   idx,
	}

	l.items[idx] = *lTask
	return true
}

func (l *listModel) SetHeight(h int) {
	l.h = h
}
func (l *listModel) SetWidth(w int) {
	l.w = w
}

func (l *listModel) SetSizes(w, h int) {
	l.SetWidth(w)
	l.SetHeight(h)
}

func (l *listModel) Counts() (int, int) {
	return l.countTotalAndCompletedTasks()
}

// WARN: private
func (l *listModel) incrementSelector() int {
	if len(l.items) < l.selectedId {
		l.selectedId += 1
	}
	return l.selectedId
}

func (l *listModel) decrementSelector() int {
	if l.selectedId > 0 {
		l.selectedId -= 1
	}
	return l.selectedId
}

// grabs up to AT_LEAST_NUMBER_OF_MONTHLY_TASKS (8) at prioritizes the uncompleted first
// IF the uncompleted tasks are equal to AT_LEAST_NUMBER_OF_MONTHLY_TASKS then it does nothing
// IF the uncompleted tasks are less to AT_LEAST_NUMBER_OF_MONTHLY_TASKS but there's no more tasks, it does nothing
// IF the uncompleted tasks are less AND there's more tasks, it just fills with whatever task there is
func (l *listModel) selectMonthlyTasksCompletedAndFill() []listItem {
	if len(l.items) == 0 {
		return []listItem{}
	}

	var atleast_monthly = map[int]listItem{}
	for _, v := range l.items {
		if v.item.Completed() {
			continue
		}
		atleast_monthly[v.item.ID()] = v

		if len(atleast_monthly) == AT_LEAST_NUMBER_OF_MONTHLY_TASKS {
			break
		}
	}

	if len(atleast_monthly) < AT_LEAST_NUMBER_OF_MONTHLY_TASKS && len(l.items) >= len(atleast_monthly) {
		for _, v := range l.items {
			_, ok := atleast_monthly[v.item.ID()]
			if ok || !v.item.Completed() {
				continue
			}
			atleast_monthly[v.item.ID()] = v

			if len(atleast_monthly) == AT_LEAST_NUMBER_OF_MONTHLY_TASKS {
				break
			}

		}
	}

	var arr_atleast_monthly = []listItem{}
	for _, v := range atleast_monthly {
		arr_atleast_monthly = append(arr_atleast_monthly, v)
	}

	//i dont know if the slices.SortedFunc is what panics since it tries to access the second item?
	if len(arr_atleast_monthly) < 2 {
		return arr_atleast_monthly
	}

	return slices.SortedFunc(slices.Values(arr_atleast_monthly), func(li1, li2 listItem) int {
		if li1.item.Completed() != li2.item.Completed() {
			if !li1.item.Completed() {
				return -1
			}
			return 1
		}

		return cmp.Compare(li1.item.ID(), li2.item.ID())
	})
}

var (
	idx_box    = lipgloss.NewStyle().Inline(true).Width(8).MaxWidth(8)
	figlet_art = lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.NewStyle().Foreground(lipgloss.Yellow).Render(FLASH_FIGLET),
		lipgloss.NewStyle().MarginTop(1).Foreground(lipgloss.BrightBlack).Render("No tasks assigned for this month"),
	)
)

const (
	max_char_until_ellipsis = 22
)

var (
	monthly_times_done = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)
)

func (_ *listModel) renderMonthlyGrid(items []listItem, c lipgloss.Style) string {
	if len(items) == 0 {
		return figlet_art
	}

	var renderedCells []string
	for _, i := range items {
		idx_cmp := idx_box.Foreground(lipgloss.Yellow)
		title_cmp := lipgloss.NewStyle().Width(max_char_until_ellipsis)

		if i.item.Completed() {
			idx_cmp = idx_cmp.Foreground(lipgloss.BrightBlack)
			title_cmp = title_cmp.Foreground(lipgloss.BrightBlack)
		}

		main_body := lipgloss.JoinHorizontal(
			lipgloss.Left,
			idx_cmp.Render(
				fmt.Sprintf("%s %d)",
					i.item.ReturnCheckboxString(), i.item.ID()),
			),
			title_cmp.Render(i.item.TitleEllipsis(max_char_until_ellipsis)),
		)

		if mTask, ok := i.item.(DZMonthlyTask); ok && mTask.TimesTotal() > 0 {
			renderedCells = append(renderedCells,
				c.Render(
					lipgloss.JoinVertical(
						lipgloss.Center,
						main_body,
						monthly_times_done.Render(fmt.Sprintf(
							"└─> %2.f / %1.f%s", mTask.CurrentProgress(), mTask.TimesTotal(), mTask.Metric()),
						),
					),
				),
			)

			continue
		}

		renderedCells = append(renderedCells,
			c.Render(main_body))

	}

	var rows []string
	for i := 0; i < len(renderedCells); i += MAX_PER_ROW {
		end := min(i+MAX_PER_ROW, len(renderedCells))
		rows = append(rows,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				renderedCells[i:end]...,
			),
		)
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (l *listModel) countTotalAndCompletedTasks() (total int, completed int) {
	completed = 0
	for _, v := range l.items {
		if v.item.Completed() {
			completed++
		}
	}

	return len(l.items), completed
}

// func (l *listModel) SelectLLTNextToExpire(){}

const (
	AT_LEAST_NUMBER_OF_MONTHLY_TASKS = 8
	MAX_PER_ROW                      = 2
)

func (m *listModel) View() string {
	monthlyItems := m.selectMonthlyTasksCompletedAndFill()
	if len(monthlyItems) == 0 {
		return ""
	}

	cWidth := m.w / 2
	cell := cStyle.Width(cWidth).MaxWidth(cWidth)
	return m.renderMonthlyGrid(monthlyItems, cell)
}
