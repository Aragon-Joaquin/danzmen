package tui

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type dzDelegate struct {
	styles styles
}

func (d dzDelegate) Height() int                             { return 1 }
func (d dzDelegate) Spacing() int                            { return 0 }
func (d dzDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

var (
	arrow_item = lipgloss.NewStyle().Foreground(lipgloss.Yellow).Padding(0, 1, 0, 2)
	id_block   = lipgloss.NewStyle().Width(5).MaxWidth(5)
)

func (d dzDelegate) Render(w io.Writer, m list.Model, idx int, listItem list.Item) {
	i, ok := listItem.(DZItem)
	if !ok {
		return
	}

	//NOTE: imagine this is a ternary
	var checked = " "
	if i.completed {
		checked = "x"
	}

	str := fmt.Sprintf("[%s] %s", checked, i.Title())
	id := id_block.Inline(true).Render(fmt.Sprintf("%d) ", i.id))

	//NOTE: select item
	if idx == m.Index() {
		style := d.styles.selectedItem
		if i.completed {
			style = style.Foreground(lipgloss.BrightBlack)
		} else {
			style = style.Underline(true)
		}

		fmt.Fprint(w, arrow_item.Render(">"), style.Render(id, str))
		return
	}

	style := d.styles.item
	if i.completed {
		style = style.Foreground(lipgloss.BrightBlack)
	}

	fmt.Fprint(w, arrow_item.Render(" "), style.Render(id, str))
}
