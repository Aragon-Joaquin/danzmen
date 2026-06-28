package tui

import (
	"fmt"
	"io"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type DzDelegateOpts struct {
	QuitImmediately bool
	RenderCompletes bool
}

type dzDelegate struct {
	styles styles
}

// i cannot put it in the struct. this is so bad
var opts *DzDelegateOpts = &DzDelegateOpts{}

func (d dzDelegate) Height() int                             { return 1 }
func (d dzDelegate) Spacing() int                            { return 0 }
func (d dzDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d dzDelegate) SetOpts(o *DzDelegateOpts) {
	if opts != nil {
		opts = o
	}
}
func (d dzDelegate) GetOpts() *DzDelegateOpts {
	return opts
}

var (
	arrow_item   = lipgloss.NewStyle().Foreground(lipgloss.Yellow).Padding(0, 1, 0, 2)
	id_block     = lipgloss.NewStyle().Width(5).MaxWidth(5).Inline(true).Foreground(lipgloss.Yellow)
	toggle_block = lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Yellow)
)

func (d dzDelegate) Render(w io.Writer, m list.Model, idx int, listItem list.Item) {
	i, ok := listItem.(DZItem)
	if !ok {
		return
	}
	//NOTE: imagine this is a ternary
	str := i.Title()
	id := id_block.Render(fmt.Sprintf("%d) ", i.id))
	toggle := toggle_block

	//NOTE: select item
	if idx == m.Index() {
		style := d.styles.selectedItem
		if i.completed {
			toggle = toggle.Foreground(lipgloss.BrightBlack)
		}

		var arrow string = ">"
		if opts.QuitImmediately {
			arrow = " "
		}
		fmt.Fprint(w, arrow_item.Render(arrow),
			style.Render(id,
				toggle.Render(i.ReturnCheckboxString()),
				str))
		return
	}

	//NOTE: normal item
	style := d.styles.item
	if i.completed {
		toggle = toggle.Foreground(lipgloss.BrightBlack)
	}

	fmt.Fprint(w, arrow_item.Render(" "),
		style.Render(
			id,
			toggle.Render(i.ReturnCheckboxString()),
			str))
}
