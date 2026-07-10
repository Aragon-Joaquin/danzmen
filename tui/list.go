package tui

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"charm.land/lipgloss/v2"
)

type DZList interface {
	GetSelectID() int
	SelectedItem() (item DZTask, ok bool)
	SetItem(idx int, item DZTask) bool

	//ui things
	SetHeight(h int)
	SetWidth(w int)
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

func CreateDZList(i []DZTask, s styles, w, h int) DZList {
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

// grabs up to AT_LEAST_NUMBER_OF_DAILY_TASKS (8) at prioritizes the uncompleted first
// IF the uncompleted tasks are equal to AT_LEAST_NUMBER_OF_DAILY_TASKS then it does nothing
// IF the uncompleted tasks are less to AT_LEAST_NUMBER_OF_DAILY_TASKS but there's no more tasks, it does nothing
// IF the uncompleted tasks are less AND there's more tasks, it just fills with whatever task there is
func (l *listModel) selectDailyTasksCompletedAndFill() []listItem {
	if len(l.items) == 0 {
		return []listItem{}
	}

	var atleast_daily = map[int]listItem{}
	for _, v := range l.items {
		if v.item.Completed() {
			continue
		}
		atleast_daily[v.item.ID()] = v

		if len(atleast_daily) == AT_LEAST_NUMBER_OF_DAILY_TASKS {
			break
		}
	}

	if len(atleast_daily) < AT_LEAST_NUMBER_OF_DAILY_TASKS && len(l.items) >= len(atleast_daily) {
		for _, v := range l.items {
			_, ok := atleast_daily[v.item.ID()]
			if ok || !v.item.Completed() {
				continue
			}
			atleast_daily[v.item.ID()] = v

			if len(atleast_daily) == AT_LEAST_NUMBER_OF_DAILY_TASKS {
				break
			}

		}
	}

	var arr_atleast_daily = []listItem{}
	for _, v := range atleast_daily {
		arr_atleast_daily = append(arr_atleast_daily, v)
	}

	//i dont know if the slices.SortedFunc is what panics since it tries to access the second item?
	if len(arr_atleast_daily) < 2 {
		return arr_atleast_daily
	}

	return slices.SortedFunc(slices.Values(arr_atleast_daily), func(li1, li2 listItem) int {
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
		lipgloss.NewStyle().MarginTop(1).Foreground(lipgloss.BrightBlack).Render("No tasks assigned for today"),
	)
)

func (_ *listModel) renderDailyGrid(items []listItem, c lipgloss.Style) string {
	if len(items) == 0 {
		return figlet_art
	}

	var renderedCells []string
	for _, i := range items {
		idx_cmp := idx_box.Foreground(lipgloss.Yellow)
		title_cmp := lipgloss.NewStyle()
		if i.item.Completed() {
			idx_cmp = idx_cmp.Foreground(lipgloss.BrightBlack)
			title_cmp = title_cmp.Foreground(lipgloss.BrightBlack)
		}

		renderedCells = append(renderedCells,
			c.Render(
				idx_cmp.Render(
					fmt.Sprintf("%s %d)",
						i.item.ReturnCheckboxString(), i.item.ID()),
				),
				title_cmp.Render(i.item.TitleEllipsis(22)),
			),
		)
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

const (
	MINIMUM_WIDTH_REQUIRED             = 80
	MINIMUM_DOUBLE_TASK_WIDTH_REQUIRED = 160
	AT_LEAST_NUMBER_OF_DAILY_TASKS     = 8
	MAX_PER_ROW                        = 2
)

var (
	baseTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Border(lipgloss.RoundedBorder(), false, false, true, false).
			AlignHorizontal(lipgloss.Center)

	dailyTitle = baseTitleStyle.
			BorderForeground(lipgloss.Yellow).
			Foreground(lipgloss.Yellow).
			MarginLeft(2)

	longTermTitle = baseTitleStyle.
			BorderForeground(lipgloss.Red).
			Foreground(lipgloss.Red).
			MarginRight(2)

	separatorLine = lipgloss.NewStyle().
			Foreground(lipgloss.BrightBlack).
			Padding(0, 2)

	cStyle = lipgloss.NewStyle().
		MarginLeft(3).
		Height(2).
		MaxHeight(2)

	remainingTasks = lipgloss.NewStyle().
			Foreground(lipgloss.BrightBlack).
			AlignHorizontal(lipgloss.Center)

	//simple UI - half the screen
	lttNotify = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Right).
			Foreground(lipgloss.BrightRed).
			MarginRight(2)

	dailyTitleHalf = dailyTitle.
			Border(lipgloss.Border{}, false).
			AlignHorizontal(lipgloss.Left).
			MarginLeft(2)

	borderBottom = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true, false, false, false).
			BorderForeground(lipgloss.Yellow)
)

// NOTE: view
func (m *listModel) View() string {
	//NOTE: not enough space
	if m.w < MINIMUM_WIDTH_REQUIRED {
		return "not enough space"
	}

	daily_itemsToRender := m.selectDailyTasksCompletedAndFill()

	total, completed := m.countTotalAndCompletedTasks()
	dailyText := fmt.Sprintf("Daily tasks (%d/%d completed)", completed, total)

	var cWidth int = 0
	var titlePadding int = 0

	//screen is bigger than 50% screen, else its smoll (<50% of screen width)
	if m.w > MINIMUM_DOUBLE_TASK_WIDTH_REQUIRED {
		cWidth = (m.w / 2) / 2 // half the width - 4 (padding) / 2 (items horizontally)
		titlePadding = (m.w - 8) / 2
	} else {
		cWidth = m.w
		titlePadding = m.w - 4
	}

	cell := cStyle.Width(cWidth).MaxWidth(cWidth)
	cellsRendered := m.renderDailyGrid(daily_itemsToRender, cell)

	var hasItemsPosition lipgloss.Position = lipgloss.Left
	if len(cellsRendered) > 0 {
		hasItemsPosition = lipgloss.Center
	}

	//NOTE: render simple UI
	if m.w < MINIMUM_DOUBLE_TASK_WIDTH_REQUIRED {
		widthForTitle := (m.w - 4) / 2 //4 for extra padding

		return lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.JoinHorizontal(
				hasItemsPosition,
				dailyTitleHalf.Width(widthForTitle).Render(dailyText),
				lttNotify.Width(widthForTitle).Render("LTT Ends in: 123d"),
			),
			borderBottom.Width(m.w).Render(),
			cellsRendered,
		)
	}

	var r_tasks string = ""
	if len(daily_itemsToRender) > AT_LEAST_NUMBER_OF_DAILY_TASKS {
		r_tasks = remainingTasks.Width((m.w - 2) / 2).Render(
			fmt.Sprintf("Show %d more tasks", len(daily_itemsToRender)-AT_LEAST_NUMBER_OF_DAILY_TASKS))
	}

	//NOTE: render complex ui (double tasks)
	dailySection := lipgloss.JoinVertical(
		hasItemsPosition,
		dailyTitle.Width(titlePadding).Render(dailyText),
		cellsRendered,
		r_tasks,
	)

	verticalBar := strings.TrimSuffix(strings.Repeat("│\n", 12), "\n")

	longTermSection := lipgloss.JoinVertical(
		lipgloss.Center,
		longTermTitle.Width(titlePadding).Render(
			fmt.Sprintf("Long term (%dd left!)", 320),
		))

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		dailySection,
		separatorLine.Render(verticalBar),
		longTermSection,
	)
}
