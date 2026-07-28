package tui

import (
	"danzmen/db"
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type TuiModel struct {
	db         *db.SqliteDB
	w          int
	h          int
	daily_list DZList
	long_list  DZList
}

const (
	DEFAULT_WIDTH = 50
	LIST_HEIGHT   = 20
)

func RenderList(daily []DZTask, long []DZLongTask, w, h int) string {
	daily_list := CreateDZList(daily, NewSimpleStyle(), w, h)
	long_list := CreateDZLongList(long, NewSimpleStyle(), w, h)
	return RenderModelView(daily_list, long_list, w, h)
}

func CreateTUIModel(daily []DZTask, long []DZLongTask, db *db.SqliteDB) TuiModel {
	mTui := TuiModel{
		db:         db,
		w:          DEFAULT_WIDTH,
		h:          LIST_HEIGHT,
		daily_list: CreateDZList(daily, NewSimpleStyle(), DEFAULT_WIDTH, LIST_HEIGHT),
		long_list:  CreateDZLongList(long, NewSimpleStyle(), DEFAULT_WIDTH, LIST_HEIGHT),
	}

	return mTui
}

func (m TuiModel) Init() tea.Cmd {
	return nil
}

func (m TuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.h = msg.Height
		m.w = msg.Width
		m.daily_list.SetSizes(msg.Width, msg.Height)
		m.long_list.SetSizes(msg.Width, msg.Height)

	case tea.KeyPressMsg:
		switch k := msg.String(); k {
		case "ctrl+c":
			return m, tea.Quit
		case "enter", "space":
			i, ok := m.daily_list.SelectedItem()

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

	v := tea.NewView(c.Render(RenderModelView(m.daily_list, m.long_list, m.w, m.h)))
	v.AltScreen = true
	return v
}
