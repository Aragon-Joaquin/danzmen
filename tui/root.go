package tui

import (
	"danzmen/db"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type tuiModel struct {
	list  list.Model
	db    *db.SqliteDB
	style styles
}

const (
	defaultWidth = 20
	listHeight   = 14
)

func CreateTUIModel(i []list.Item, db *db.SqliteDB, s styles) tuiModel {
	l := list.New(i, dzDelegate{&s}, defaultWidth, listHeight)

	l.SetShowTitle(true)
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(false)

	l.Title = "Tasks pending today"
	l.InfiniteScrolling = true
	l.Styles.Title = lipgloss.NewStyle()

	mTui := tuiModel{list: l, style: s, db: db}
	mTui.updateStyles(s)
	return mTui
}

func (m tuiModel) Init() tea.Cmd {
	m.list.Items()
	return nil
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
	case tea.KeyPressMsg:

		switch k := msg.String(); k {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			_, ok := m.list.SelectedItem().(DZItem)

			if !ok {
				break
			}

		}
	}

	var cmd tea.BatchMsg
	l, c := m.list.Update(msg)
	m.list = l
	cmd = append(cmd, c)
	return m, tea.Batch(cmd...)
}

func (m tuiModel) View() tea.View {
	v := tea.NewView(m.list.View())
	v.AltScreen = true
	return v
}

// NOTE: private
func (m tuiModel) updateStyles(s styles) {
	m.style = s
	m.list.Styles.Title = m.style.title
	m.list.Styles.PaginationStyle = m.style.pagination
	m.list.Styles.HelpStyle = m.style.help
	m.list.SetDelegate(dzDelegate{styles: &m.style})
}
