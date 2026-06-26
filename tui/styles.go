package tui

import (
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"
)

type styles struct {
	title        lipgloss.Style
	item         lipgloss.Style
	selectedItem lipgloss.Style
	pagination   lipgloss.Style
	help         lipgloss.Style
	quitText     lipgloss.Style
}

func NewSimpleStyle() styles {
	var s styles
	s.title = lipgloss.NewStyle().Margin(1, 2, 0, 0).Padding(0, 1).Foreground(lipgloss.Yellow).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Yellow)
	s.item = lipgloss.NewStyle()
	s.selectedItem = lipgloss.NewStyle().Foreground(lipgloss.Yellow)
	s.pagination = list.DefaultStyles(false).PaginationStyle.PaddingLeft(4)
	s.help = list.DefaultStyles(false).HelpStyle.PaddingLeft(4).PaddingBottom(1)
	s.quitText = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	return s
}
