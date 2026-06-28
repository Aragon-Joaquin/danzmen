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
	s.selectedItem = lipgloss.NewStyle()
	s.pagination = list.DefaultStyles(false).PaginationStyle.PaddingLeft(4)
	s.help = list.DefaultStyles(false).HelpStyle.PaddingLeft(4).PaddingBottom(1)
	s.quitText = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	return s
}

func CreateListModel(i []list.Item, s styles, o *DzDelegateOpts) list.Model {
	dz := dzDelegate{styles: s}
	dz.SetOpts(o)

	l := list.New(i, dz, DEFAULT_WIDTH, LIST_HEIGHT)

	l.SetShowTitle(false)
	//l.Title = "Tasks pending today"

	l.SetFilteringEnabled(!dz.GetOpts().QuitImmediately)
	l.SetShowStatusBar(false)

	//TODO: make my own
	l.SetShowHelp(!dz.GetOpts().QuitImmediately)

	l.InfiniteScrolling = true
	return updateListStyles(&l, s)
}

// NOTE: private
func updateListStyles(list *list.Model, s styles) list.Model {
	l := list
	l.Styles.NoItems = l.Styles.NoItems.
		MarginLeft(4)

	l.Styles.Title = s.title
	l.Styles.PaginationStyle = s.pagination
	l.Styles.HelpStyle = s.help
	l.SetDelegate(dzDelegate{styles: s})
	return *l
}
