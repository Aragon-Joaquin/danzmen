package tui

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type dzDelegate struct {
	styles *styles
}

func (d dzDelegate) Height() int                             { return 1 }
func (d dzDelegate) Spacing() int                            { return 0 }
func (d dzDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d dzDelegate) Render(w io.Writer, m list.Model, idx int, listItem list.Item) {
	i, ok := listItem.(DZItem)
	if !ok {
		return
	}

	var checked = " "
	if i.completed {
		checked = "x"
	}

	str := fmt.Sprintf("%d) [%s] %s", i.id, checked, i.Title())

	fn := d.styles.item.Render
	if idx == m.Index() {
		fn = func(s ...string) string {
			return d.styles.selectedItem.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
