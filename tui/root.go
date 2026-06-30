package tui

import (
	"danzmen/db"
	"log"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type TuiModel struct {
	list  list.Model
	db    *db.SqliteDB
	style styles

	quitImmediately bool

	w int
	h int
}

const (
	DEFAULT_WIDTH = 50
	LIST_HEIGHT   = 8
)

func CreateTUIModel(i []list.Item, db *db.SqliteDB, s styles, q bool) TuiModel {
	d := &DzDelegateOpts{QuitImmediately: q, RenderCompletes: !q}
	mTui := TuiModel{
		list:            CreateListModel(i, s, d),
		style:           s,
		db:              db,
		quitImmediately: q,
	}

	return mTui
}

func (m TuiModel) Init() tea.Cmd {
	if m.quitImmediately {
		return tea.Quit
	}
	return nil
}

func (m TuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.h = msg.Height
		m.w = msg.Width
		m.list.SetWidth(msg.Width)

	case tea.KeyPressMsg:
		switch k := msg.String(); k {
		case "ctrl+c":
			return m, tea.Quit
		case "enter", "space":
			i, ok := m.list.SelectedItem().(DZItem)

			if !ok {
				break
			}

			//uncheck by default
			var toggle int = 0
			if !i.completed {
				toggle = 1
			}

			if err := m.db.UpdateCompletedTask(i.id, toggle); err != nil {
				log.Println(err)
				return m, nil
			}

			i.completed = !i.completed
			m.list.SetItem(m.list.Index(), i)
		}
	}

	var cmd tea.BatchMsg
	l, c := m.list.Update(msg)
	m.list = l
	cmd = append(cmd, c)
	return m, tea.Batch(cmd...)
}

var (
	container = lipgloss.
		NewStyle().
		Height(LIST_HEIGHT).
		MaxHeight(LIST_HEIGHT)
)

func (m TuiModel) View() tea.View {
	c := container.Width(m.w).MarginTop(1).Padding(0)

	v := tea.NewView(c.Render(m.list.View()))
	if !m.quitImmediately {
		v.AltScreen = true
	}
	return v
}
