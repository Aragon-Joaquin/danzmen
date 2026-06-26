package tui

import (
	"danzmen/db"
	"log"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type tuiModel struct {
	list  list.Model
	db    *db.SqliteDB
	style styles

	w int
	h int
}

const (
	defaultWidth = 20
	listHeight   = 14
)

func CreateTUIModel(i []list.Item, db *db.SqliteDB, s styles) tuiModel {
	l := list.New(i, dzDelegate{s}, defaultWidth, listHeight)

	l.SetShowTitle(true)
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(false)

	//TODO: make my own
	l.SetShowHelp(false)

	l.Title = "Tasks pending today"
	l.InfiniteScrolling = true

	mTui := tuiModel{list: l, style: s, db: db}
	mTui.updateStyles(s)

	return mTui
}

func (m tuiModel) Init() tea.Cmd {
	return nil
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.h = msg.Height
		m.w = msg.Width
		m.list.SetWidth(msg.Width)

	case tea.KeyPressMsg:
		switch k := msg.String(); k {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(DZItem)

			if !ok {
				break
			}

			var toggle db.DBComplete = db.Completed_NO
			if i.completed {
				toggle = db.Completed_YES
			}

			err := m.db.UpdateCompletedTask(i.id, toggle)
			if err != nil {
				log.Println(err)
				return m, nil
			}
			i.completed = !toggle.IsCompleted()

		}
	}

	var cmd tea.BatchMsg
	l, c := m.list.Update(msg)
	m.list = l
	cmd = append(cmd, c)
	return m, tea.Batch(cmd...)
}

const (
	MAX_HEIGHT_CONTAINER = 20
)

var (
	container = lipgloss.
			NewStyle().
			Height(MAX_HEIGHT_CONTAINER).
			MaxHeight(MAX_HEIGHT_CONTAINER)

	tasks_completed = lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Height(MAX_HEIGHT_CONTAINER).
			Padding(2, 5)
)

func (m tuiModel) View() tea.View {
	c := container.Width(m.w)

	var content string = m.list.View()
	if len(m.list.Items()) == 0 {
		c = c.Width(m.w).Foreground(lipgloss.Yellow)
		content = tasks_completed.
			Width(m.w).
			Render("Today's tasks were completed" + "\n" + SUN_FIGLET)
	}

	v := tea.NewView(
		c.Render(content),
	)
	v.AltScreen = true
	return v
}

// NOTE: private
func (m *tuiModel) updateStyles(s styles) {
	m.style = s
	m.list.Styles.Title = s.title
	m.list.Styles.PaginationStyle = s.pagination
	m.list.Styles.HelpStyle = s.help

	m.list.SetDelegate(dzDelegate{styles: s})
}
