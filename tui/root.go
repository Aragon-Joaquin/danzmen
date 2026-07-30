package tui

import (
	"danzmen/db"
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type TuiModel struct {
	db           *db.SqliteDB
	w            int
	h            int
	monthly_list DZList
	long_list    DZList
}

const (
	DEFAULT_WIDTH = 50
	LIST_HEIGHT   = 20
)

func RenderList(monthly []DZMonthlyTask, long []DZLongTask, w, h int) string {
	monthly_list := CreateDZList(monthly, NewSimpleStyle(), w, h)
	long_list := CreateDZLongList(long, NewSimpleStyle(), w, h)
	return RenderModelView(monthly_list, long_list, w, h)
}

func CreateTUIModel(monthly []DZMonthlyTask, long []DZLongTask, db *db.SqliteDB) TuiModel {
	mTui := TuiModel{
		db:           db,
		w:            DEFAULT_WIDTH,
		h:            LIST_HEIGHT,
		monthly_list: CreateDZList(monthly, NewSimpleStyle(), DEFAULT_WIDTH, LIST_HEIGHT),
		long_list:    CreateDZLongList(long, NewSimpleStyle(), DEFAULT_WIDTH, LIST_HEIGHT),
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
		m.monthly_list.SetSizes(msg.Width, msg.Height)
		m.long_list.SetSizes(msg.Width, msg.Height)

	case tea.KeyPressMsg:
		switch k := msg.String(); k {
		case "ctrl+c":
			return m, tea.Quit
		case "enter", "space":
			i, ok := m.monthly_list.SelectedItem()

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
	CONTAINER_FOR_TOGGLE = lipgloss.
				NewStyle().
				Margin(1, 0).
				Padding(0)

	c = CONTAINER_FOR_TOGGLE.
		Height(LIST_HEIGHT).
		MaxHeight(LIST_HEIGHT)
)

func (m TuiModel) View() tea.View {
	v := tea.NewView(
		c.Render(RenderModelView(m.monthly_list, m.long_list, m.w, m.h)))
	v.AltScreen = true
	return v
}
