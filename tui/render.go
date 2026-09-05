package tui

import (
	"danzmen/db"
	ty "danzmen/types"
	"strings"

	"charm.land/lipgloss/v2"
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
			Padding(0, 1)

	cStyle = lipgloss.NewStyle().
		Height(3).
		MaxHeight(3).
		Align(lipgloss.Left)

	remainingTasks = lipgloss.NewStyle().
			Foreground(lipgloss.BrightBlack).
			AlignHorizontal(lipgloss.Center)

	lttTitle    = lipgloss.NewStyle()
	lttidxbox   = idx_box.Foreground(lipgloss.BrightRed)
	lttPriority = lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true).
			MaxWidth(6).
			Width(6).
			Align(lipgloss.Center)
	lttEnds = lipgloss.NewStyle().Padding(0, 1)

	//simple UI - half the screen
	lttNotify = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Right).
			Foreground(lipgloss.BrightRed).
			MarginRight(2)

	monthlyTitleHalf = monthlyTitle.
				Border(lipgloss.Border{}, false).
				AlignHorizontal(lipgloss.Left)

	borderBottom = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true, false, false, false).
			BorderForeground(lipgloss.Yellow)

	ltt_empty_task_placeholder = "Nothing here !"
)

func RenderModelView(monthlyI DZList, longI DZList, w, h int, displayData *db.DisplayData) string {
	half, opts := GetTypeOfSize(w)
	//NOTE: nothing
	if opts.typ == SIZE_SMALL {
		return ""
	}

	ll, ml, err := TypeCastToListModel(longI, monthlyI)
	if err != nil {
		return err.Error()
	}

	monthlyItems := ml.selectTasksCompletedAndFill(ty.AT_LEAST_NUMBER_OF_MONTHLY_TASKS)
	monthlyText := GenerateMonthlyText(ml, displayData.CountMonthlyTasks)

	cellsRendered := ml.renderMonthlyGrid(monthlyItems,
		cStyle.Margin(0, 1).Width(opts.cellWidth).MaxWidth(opts.cellWidth))

	var hasItemsPosition lipgloss.Position = lipgloss.Left
	if len(cellsRendered) > 0 {
		hasItemsPosition = lipgloss.Center
	}

	nt_placeholder := ll.nearest_task_placeholder(displayData.NearestLongTaskToExpire)

	switch opts.typ {
	//NOTE: daily
	case SIZE_MEDIUM:
		ml.SetWidth(w)

		return lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.JoinHorizontal(
				hasItemsPosition,
				monthlyTitleHalf.Width(opts.titlePadding).Render(monthlyText),
				lttNotify.Width(opts.titlePadding).Render(nt_placeholder),
			),
			borderBottom.Width(w-1).Render(),
			cellsRendered,
		)

		//NOTE: daily | long
	case SIZE_BIG:
		ml.SetWidth(half)

		//shows if there's more monthly tasks remain to be rendered
		r_tasks := RenderPlaceholderTextShowingRemainingTasksToRender(w, len(monthlyItems))

		monthlySection := lipgloss.JoinVertical(
			hasItemsPosition,
			monthlyTitle.Width(opts.titlePadding).Render(monthlyText),
			cellsRendered,
			r_tasks,
		)

		verticalBar := strings.TrimSuffix(strings.Repeat("│\n", 14), "\n")
		longTermSection := opts.ReturnLongListRenderable(ll, displayData.NearestLongTaskToExpire)

		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			lipgloss.NewStyle().Render(monthlySection),
			separatorLine.Render(verticalBar), //vertical bar
			lipgloss.NewStyle().Render(longTermSection),
		)

	default:
		return "entered unexpected case"
	}
}

// this is more understandable than the debugger output
// return lipgloss.JoinHorizontal(
// 	lipgloss.Left,
// 	lipgloss.NewStyle().Background(lipgloss.Blue).Render(monthlySection),
// 	separatorLine.Background(lipgloss.Green).Render(verticalBar),
// 	lipgloss.NewStyle().Background(lipgloss.Red).Render(longTermSection),
// )
