package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sno6/brain/search"
)

// Page is an enum type that we use to determine what components to render.
// This type can be sent as a tea.Msg from sub-components to change the page.
type Page uint8

const (
	PageIndex Page = iota
	PageSearch
	PageWrite
	PageView
)

func changePage(p Page) func() tea.Msg {
	return func() tea.Msg {
		return p
	}
}

// A savedCell is a message type that is passed to an App update
// when the user saves a cell.
type savedCell string

func saveCell(c savedCell) func() tea.Msg {
	return func() tea.Msg {
		return c
	}
}

// A searchMessage contains the contents of the search bar, and is
// sent to other models when the user stops typing briefly.
type searchMessage struct {
	mode search.Mode
	val  string
}

func searchCommand(mode search.Mode, val string) func() tea.Msg {
	return func() tea.Msg {
		return searchMessage{
			mode: mode,
			val:  val,
		}
	}
}

type viewCellMessage string

func viewCellCommand(content string) func() tea.Msg {
	return func() tea.Msg {
		return viewCellMessage(content)
	}
}
