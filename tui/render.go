package tui

import (
	"fmt"
	"strings"
	"time"

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

	monthlyTitle = baseTitleStyle.
			BorderForeground(lipgloss.Yellow).
			Foreground(lipgloss.Yellow).MarginLeft(2)

	longTermTitle = baseTitleStyle.
			BorderForeground(lipgloss.Red).
			Foreground(lipgloss.Red).
			MarginRight(2)

	separatorLine = lipgloss.NewStyle().
			Foreground(lipgloss.BrightBlack).
			Padding(0, 2)

	cStyle = lipgloss.NewStyle().
		Margin(0, 1).
		Height(2).
		MaxHeight(2)

	remainingTasks = lipgloss.NewStyle().
			Foreground(lipgloss.BrightBlack).
			AlignHorizontal(lipgloss.Center)

	lttTitle    = lipgloss.NewStyle()
	lttidxbox   = idx_box.Foreground(lipgloss.BrightRed)
	lttPriority = lipgloss.NewStyle().Padding(0, 1).Margin(0, 1)
	lttEnds     = lipgloss.NewStyle()

	//simple UI - half the screen
	lttNotify = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Right).
			Foreground(lipgloss.BrightRed).
			MarginRight(2)

	monthlyTitleHalf = monthlyTitle.
			Border(lipgloss.Border{}, false).
			AlignHorizontal(lipgloss.Left).
			MarginLeft(2)

	borderBottom = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true, false, false, false).
			BorderForeground(lipgloss.Yellow)

	ltt_empty_task_placeholder = "Nothing to worry about"
)

func RenderModelView(monthlyI DZList, longI DZList, w, h int) string {
	if w < MINIMUM_WIDTH_REQUIRED {
		return "not enough space"
	}

	var ll *listModel = &listModel{}
	if lm, ok := longI.(*listModel); ok {
		ll = lm
	}

	var ml *listModel = &listModel{}
	if lm, ok := monthlyI.(*listModel); ok {
		ml = lm
	}

	monthlyItems := ml.selectMonthlyTasksCompletedAndFill()
	total, completed := ml.countTotalAndCompletedTasks()

	monthlyText := fmt.Sprintf("Monthly tasks (%d/%d completed)", completed, total)

	var cWidth int
	var titlePadding int

	//screen is bigger than 50% screen, else its small (<50% of screen width)
	cWidth = (w / 2) - 2 // 2 for padding

	if w > MINIMUM_DOUBLE_TASK_WIDTH_REQUIRED {
		cWidth = cWidth / 2
		titlePadding = (w - 8) / 2
		ml.SetWidth(w / 2)
	} else {
		titlePadding = w - 4
		ml.SetWidth(w)
	}

	cell := cStyle.Width(cWidth).MaxWidth(cWidth)
	cellsRendered := ml.renderMonthlyGrid(monthlyItems, cell)

	var hasItemsPosition lipgloss.Position = lipgloss.Left
	if len(cellsRendered) > 0 {
		hasItemsPosition = lipgloss.Center
	}

	_, dd, ok := findLongTaskWithLeastDaysOfCompletion(ll)
	//NOTE: render simple ui
	if w < MINIMUM_DOUBLE_TASK_WIDTH_REQUIRED {
		widthForTitle := (w - 4) / 2

		var nearest_task_placeholder = ltt_empty_task_placeholder
		if ok {
			nearest_task_placeholder = fmt.Sprintf("Next LTT ends in: %dd", dd)
		}

		return lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.JoinHorizontal(
				hasItemsPosition,
				monthlyTitleHalf.Width(widthForTitle).Render(monthlyText),
				lttNotify.Width(widthForTitle).Render(nearest_task_placeholder),
			),
			borderBottom.Width(w-1).Render(),
			cellsRendered,
		)
	}

	var r_tasks string
	if len(monthlyItems) > AT_LEAST_NUMBER_OF_MONTHLY_TASKS {
		r_tasks = remainingTasks.Width((w - 2) / 2).Render(
			fmt.Sprintf("Show %d more tasks", len(monthlyItems)-AT_LEAST_NUMBER_OF_MONTHLY_TASKS))
	}

	//NOTE: render complex ui (double tasks)
	monthlySection := lipgloss.JoinVertical(
		hasItemsPosition,
		monthlyTitle.Width(titlePadding).Render(monthlyText),
		cellsRendered,
		r_tasks,
	)

	verticalBar := strings.TrimSuffix(strings.Repeat("│\n", 14), "\n")

	var longRows []string
	for _, li := range ll.items {
		lt, ok := li.item.(DZLongTask)
		if !ok {
			continue
		}
		icon := lt.ReturnCheckboxString()

		row := cStyle.Height(1).Width((cWidth*2)-1).Render(
			lttidxbox.Render(fmt.Sprintf("%s %d)", icon, lt.ID())),
			lttTitle.Width(int(float64(cWidth)*1.4)).Render(lt.TitleEllipsis(22)),
			lttPriority.Background(lt.PriorityBGColor()).Render(lt.RenderPriority()),
			lttEnds.Render(lt.HumanReadableEndsIn()),
		)

		longRows = append(longRows, row)
	}

	longContent := strings.Join(longRows, "\n")
	if longContent == "" {
		longContent = lipgloss.NewStyle().Render(ltt_empty_task_placeholder)
	}

	longTermSection := lipgloss.JoinVertical(
		lipgloss.Center,
		longTermTitle.Width(titlePadding).Render(
			fmt.Sprintf("Long term (%dd left!)", dd),
		),
		lipgloss.NewStyle().Render(longContent),
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		monthlySection,
		separatorLine.Render(verticalBar),
		longTermSection,
	)
}

// WARN: PRIVATE
// find the task with the least days to be completed
func findLongTaskWithLeastDaysOfCompletion(ll *listModel) (selectedTask DZLongTask, days_diff int, ok bool) {
	var nearestDate time.Time
	now := time.Now()

	selectedTask = nil
	ok = false
	days_diff = -1

	for _, li := range ll.items {
		lt, ok := li.item.(DZLongTask)
		if !ok {
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

	return
}
