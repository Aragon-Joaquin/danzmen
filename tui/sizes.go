package tui

import (
	ty "danzmen/types"
	"fmt"

	"charm.land/lipgloss/v2"
)

type CLI_SIZE int

const (
	SIZE_SMALL CLI_SIZE = iota
	SIZE_MEDIUM
	SIZE_BIG
)

const (
	MINIMUM_WIDTH_REQUIRED             = 80
	MINIMUM_DOUBLE_TASK_WIDTH_REQUIRED = 160
)

type SizeOpts struct {
	typ          CLI_SIZE
	cellWidth    int
	titlePadding int
}

func GetTypeOfSize(w int) (halfScreen int, opts *SizeOpts) {
	if w <= MINIMUM_WIDTH_REQUIRED {
		return 0, &SizeOpts{
			typ: SIZE_SMALL,
		}
	}

	cWidth := (w / 2) - 2 // 2 for padding

	if w < MINIMUM_DOUBLE_TASK_WIDTH_REQUIRED {
		return w, &SizeOpts{
			typ:          SIZE_MEDIUM,
			cellWidth:    cWidth,
			titlePadding: (w - 4) / 2,
		}
	}

	half := (w - 3) / 2 // vertical bar is 3 wide (1 char + padding)
	return half, &SizeOpts{
		typ:          SIZE_BIG,
		cellWidth:    (half - 4) / 2,
		titlePadding: half - 2,
	}
}

func (s *SizeOpts) ReturnLongListRenderable(ll *listModel, dd int) string {
	longContent := ll.renderLongGrid(
		ll.selectTasksCompletedAndFill(ty.AT_LEAST_NUMBER_OF_LONG_TASKS),
		cStyle.MaxHeight(2).Height(1).MarginBottom(1).Width(((s.cellWidth) * 2)), s.cellWidth)

	longTermSection := lipgloss.JoinVertical(
		lipgloss.Center,
		longTermTitle.Width(s.titlePadding-1).Render(
			fmt.Sprintf("Long term (%dd left!)", dd),
		),
		lipgloss.NewStyle().Render(longContent),
	)

	return longTermSection
}
