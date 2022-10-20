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

type deleteCellMessage string

func deleteCell(id string) func() tea.Msg {
	return func() tea.Msg {
		return deleteCellMessage(id)
	}
}

// A savedCell is a message type that is passed to an App update
// when the user saves a cell. If docID is present the user is editing
// the document and we should remove the original before writing the new value.
type savedCell struct {
	docID   string
	content string
}

func saveCell(docID, content string) func() tea.Msg {
	return func() tea.Msg {
		return savedCell{docID: docID, content: content}
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

type viewCellMessage struct {
	id, content string
}

func viewCellCommand(id, content string) func() tea.Msg {
	return func() tea.Msg {
		return viewCellMessage{id: id, content: content}
	}
}
