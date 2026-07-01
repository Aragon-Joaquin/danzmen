package tui

import "charm.land/lipgloss/v2"

type styles struct {
	//	title        lipgloss.Style
	item         lipgloss.Style
	selectedItem lipgloss.Style
	pagination   lipgloss.Style
	help         lipgloss.Style
}

func NewSimpleStyle() styles {
	var s styles
	//s.title = lipgloss.NewStyle().Margin(1, 2, 0, 0).Padding(0, 1).Foreground(lipgloss.Yellow).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Yellow)
	s.item = lipgloss.NewStyle()
	s.selectedItem = lipgloss.NewStyle()
	s.pagination = lipgloss.NewStyle().PaddingLeft(4)
	s.help = lipgloss.NewStyle().PaddingLeft(4).PaddingBottom(1)
	return s
}

