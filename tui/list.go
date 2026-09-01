package tui

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"time"

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
	// View() string
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

func (l *listModel) Counts() (int, int) { return l.countTotalAndCompletedTasks() }
func (l *listModel) SetHeight(h int)    { l.h = h }
func (l *listModel) SetWidth(w int)     { l.w = w }

func (l *listModel) SetSizes(w, h int) {
	l.SetWidth(w)
	l.SetHeight(h)
}

// func (m *listModel) View() string {
// 	monthlyItems := m.selectTasksCompletedAndFill()
// 	if len(monthlyItems) == 0 {
// 		return ""
// 	}
//
// 	cWidth := m.w / 2
// 	cell := cStyle.Width(cWidth).MaxWidth(cWidth)
// 	return m.renderMonthlyGrid(monthlyItems, cell)
// }

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

const (
	max_char_until_ellipsis = 30
	MAX_PER_ROW             = 2
)

type AT_LEAST_NUMBER int

const (
	AT_LEAST_NUMBER_OF_MONTHLY_TASKS AT_LEAST_NUMBER = 8
	AT_LEAST_NUMBER_OF_LONG_TASKS    AT_LEAST_NUMBER = 6
)

// grabs up to AT_LEAST_NUMBER_OF_MONTHLY_TASKS (8) at prioritizes the uncompleted first
// IF the uncompleted tasks are equal to AT_LEAST_NUMBER_OF_MONTHLY_TASKS then it does nothing
// IF the uncompleted tasks are less to AT_LEAST_NUMBER_OF_MONTHLY_TASKS but there's no more tasks, it does nothing
// IF the uncompleted tasks are less AND there's more tasks, it just fills with whatever task there is
func (l *listModel) selectTasksCompletedAndFill(at_least_tasks AT_LEAST_NUMBER) []listItem {
	if len(l.items) == 0 {
		return []listItem{}
	}

	AT_LEAST := int(at_least_tasks)

	var atleast = map[int]listItem{}
	for _, v := range l.items {
		if v.item.Completed() {
			continue
		}

		atleast[v.item.ID()] = v

		if len(atleast) == AT_LEAST {
			break
		}
	}

	if len(atleast) < AT_LEAST && len(l.items) >= len(atleast) {
		for _, v := range l.items {
			_, ok := atleast[v.item.ID()]
			if ok || !v.item.Completed() {
				continue
			}
			atleast[v.item.ID()] = v

			if len(atleast) == AT_LEAST {
				break
			}

		}
	}

	var arr_atleast = []listItem{}
	for _, v := range atleast {
		arr_atleast = append(arr_atleast, v)
	}

	//i dont know if the slices.SortedFunc is what panics since it tries to access the second item?
	if len(arr_atleast) < 2 {
		return arr_atleast
	}

	return slices.SortedFunc(slices.Values(arr_atleast), func(li1, li2 listItem) int {
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

var (
	monthly_times_done = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)
)

// render the monthly tasks
func (l *listModel) renderMonthlyGrid(items []listItem, c lipgloss.Style) string {
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
			title_cmp.Strikethrough(i.item.Completed()).Render(i.item.TitleEllipsis(max_char_until_ellipsis)),
		)

		if mTask, ok := i.item.(DZMonthlyTask); ok && mTask.TimesTotal() > 0 {
			renderedCells = append(renderedCells,
				c.Render(
					lipgloss.JoinVertical(
						lipgloss.Center,
						main_body,
						monthly_times_done.MarginLeft(4).Width(max_char_until_ellipsis+4).Strikethrough(i.item.Completed()).Render(fmt.Sprintf(
							"└─>%2.f / %1.f%s", mTask.CurrentProgress(), mTask.TimesTotal(), mTask.Metric()),
						),
					),
				),
			)

			continue
		}

		renderedCells = append(renderedCells,
			c.Render(main_body))

	}

	return lipgloss.JoinVertical(lipgloss.Left, (l.arrange_cells_in_layout(renderedCells, MAX_PER_ROW))...)
}

func (_ *listModel) renderLongGrid(items []listItem, c lipgloss.Style, cWidth int) string {
	var longRows []string
	for _, li := range items {
		lt, ok := li.item.(DZLongTask)
		if !ok {
			continue
		}

		icon := lt.ReturnCheckboxString()
		row := c.Render(
			lttidxbox.Render(fmt.Sprintf("%s %d)", icon, lt.ID())),
			lttTitle.Width(int(float64(cWidth)*1.3)).Render(lt.TitleEllipsis(max_char_until_ellipsis)),
			lttPriority.Background(lt.PriorityBGColor()).Render(lt.RenderPriority()),
			lttEnds.Render(lt.HumanReadableEndsIn()),
		)

		longRows = append(longRows, row)
	}

	longContent := strings.Join(longRows, "\n")
	if longContent == "" {
		longContent = lipgloss.NewStyle().Render(ltt_empty_task_placeholder)
	}

	return longContent
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

func (l *listModel) arrange_cells_in_layout(cells []string, MAX_PER_ROW int) []string {
	rows := []string{}
	for i := 0; i < len(cells); i += MAX_PER_ROW {
		end := min(i+MAX_PER_ROW, len(cells))
		rows = append(rows,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				cells[i:end]...,
			),
		)
	}

	return rows
}

// find the task with the least days to be completed
func (ll *listModel) findLongTaskNextToExpire() (selectedTask DZLongTask, days_diff int, placeholder string) {
	var nearestDate time.Time
	now := time.Now()

	var ok bool = false
	selectedTask = nil
	days_diff = -1
	placeholder = ltt_empty_task_placeholder

	for _, li := range ll.items {
		lt, k := li.item.(DZLongTask)
		if !k {
			continue
		}

		t, err := lt.MM_DD_YYYY_Format()
		if err != nil {
			continue
		}

		if t.After(now) && (nearestDate.IsZero() || t.Before(nearestDate)) {
			selectedTask = lt
			ok = true
			days_diff = int(t.Sub(now).Hours() / 24)
		}
	}

	if ok {
		placeholder = fmt.Sprintf("Next LTT ends in: %dd", days_diff)
	}

	return
}
