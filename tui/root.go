package tui

import (
	"danzmen/db"
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type TuiModel struct {
	db              *db.SqliteDB
	quitImmediately bool
	w               int
	h               int
	list            DZList
}

const (
	DEFAULT_WIDTH = 50
	LIST_HEIGHT   = 20
)

func CreateTUIModel(i []DZTask, db *db.SqliteDB, q bool) TuiModel {
	mTui := TuiModel{
		db:              db,
		quitImmediately: q,
		w:               DEFAULT_WIDTH,
		h:               LIST_HEIGHT,
		list:            CreateDZList(i, NewSimpleStyle(), DEFAULT_WIDTH, LIST_HEIGHT),
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
		m.list.SetHeight(msg.Height)
		m.list.SetWidth(msg.Width)

	case tea.KeyPressMsg:
		switch k := msg.String(); k {
		case "ctrl+c":
			return m, tea.Quit
		case "enter", "space":
			i, ok := m.list.SelectedItem()

			if !ok {
				break
			}

			var toggle int = 0
			if !i.Completed() {
				toggle = 1
			}

			if err := m.db.UpdateCompletedTask(i.ID(), toggle); err != nil {
				log.Println(err)
				return m, nil
			}

			//i.completed = !i.completed
		}
	}

	var cmd tea.BatchMsg
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
