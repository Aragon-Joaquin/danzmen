package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

const (
	MINIMUM_WIDTH_REQUIRED             = 80
	MINIMUM_DOUBLE_TASK_WIDTH_REQUIRED = 160
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

func RenderModelView(dailyI DZList, longI DZList, w, h int) string {
	if w < MINIMUM_WIDTH_REQUIRED {
		return "not enough space"
	}

	var ll *listModel = &listModel{}
	if lm, ok := longI.(*listModel); ok {
		longI = lm
	}

	var dl *listModel = &listModel{}
	if lm, ok := dailyI.(*listModel); ok {
		dailyI = lm
	}

	dailyItems := dl.selectDailyTasksCompletedAndFill()
	total, completed := dl.countTotalAndCompletedTasks()

	dailyText := fmt.Sprintf("Daily tasks (%d/%d completed)", completed, total)

	var cWidth int
	var titlePadding int

	//screen is bigger than 50% screen, else its smoll (<50% of screen width)
	cWidth = w / 2

	if w > MINIMUM_DOUBLE_TASK_WIDTH_REQUIRED {
		cWidth = cWidth / 2
		titlePadding = (w - 8) / 2
		dl.SetWidth(w / 2)
	} else {
		titlePadding = w - 4
		dl.SetWidth(w)
	}

	cell := cStyle.Width(cWidth).MaxWidth(cWidth)
	cellsRendered := dl.renderDailyGrid(dailyItems, cell)

	var hasItemsPosition lipgloss.Position = lipgloss.Left
	if len(cellsRendered) > 0 {
		hasItemsPosition = lipgloss.Center
	}

	//NOTE: render simple ui
	if w < MINIMUM_DOUBLE_TASK_WIDTH_REQUIRED {
		widthForTitle := (w - 4) / 2

		return lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.JoinHorizontal(
				hasItemsPosition,
				dailyTitleHalf.Width(widthForTitle).Render(dailyText),
				lttNotify.Width(widthForTitle).Render(fmt.Sprintf("Next LTT ends in: %s", "lol")),
			),
			borderBottom.Width(w-1).Render(),
			cellsRendered,
		)
	}

	var r_tasks string
	if len(dailyItems) > AT_LEAST_NUMBER_OF_DAILY_TASKS {
		r_tasks = remainingTasks.Width((w - 2) / 2).Render(
			fmt.Sprintf("Show %d more tasks", len(dailyItems)-AT_LEAST_NUMBER_OF_DAILY_TASKS))
	}

	//NOTE: render complex ui (double tasks)
	dailySection := lipgloss.JoinVertical(
		hasItemsPosition,
		dailyTitle.Width(titlePadding).Render(dailyText),
		cellsRendered,
		r_tasks,
	)

	verticalBar := strings.TrimSuffix(strings.Repeat("│\n", 10), "\n")

	var longRows []string
	for _, li := range ll.items {
		if lt, ok := li.item.(DZLongTask); ok {
			icon := "[ ]"
			if lt.Completed() {
				icon = "[x]"
			}
			row := fmt.Sprintf("%s %s (ends: %s)", icon, lt.TitleEllipsis(22), lt.EndsOn())
			longRows = append(longRows, row)
		} else {
			icon := "[ ]"
			if li.item.Completed() {
				icon = "[x]"
			}
			row := fmt.Sprintf("%s %s", icon, li.item.TitleEllipsis(22))
			longRows = append(longRows, row)
		}
	}

	longContent := strings.Join(longRows, "\n")
	if longContent == "" {
		longContent = "No long-term tasks"
	}

	longTermSection := lipgloss.JoinVertical(
		lipgloss.Center,
		longTermTitle.Width(titlePadding).Render(
			fmt.Sprintf("Long term (%dd left!)", 320),
		),
		lipgloss.NewStyle().PaddingLeft(2).Render(longContent),
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		dailySection,
		separatorLine.Render(verticalBar),
		longTermSection,
	)
}
