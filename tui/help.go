package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type helpModel struct {
	keyMap keyMap
	help   help.Model
}

func newHelpModel(page Page) *helpModel {
	return &helpModel{
		keyMap: buildKeyMap(page),
		help:   help.New(),
	}
}

func (h *helpModel) Init() tea.Cmd {
	return nil
}

func (h *helpModel) Update(msg tea.Msg) (*helpModel, tea.Cmd) {
	return h, nil
}

func (h *helpModel) View() string {
	if h.keyMap.page == PageSearch {
		return h.help.View(h.keyMap)
	}
	return " " + h.help.View(h.keyMap)
}

func buildKeyMap(page Page) keyMap {
	return keyMap{
		page: page,
		Save: key.NewBinding(
			key.WithKeys("ctrl+s", "ctrl+s"),
			key.WithHelp("ctrl+s", "save"),
		),
		ToggleSearch: key.NewBinding(
			key.WithKeys("tab", "tab"),
			key.WithHelp("tab", "toggle search"),
		),
		// Close a single view - stay in app.
		Quit: key.NewBinding(
			key.WithKeys("q", "q"),
			key.WithHelp("q", "close view"),
		),
		// Exit the whole app.
		Exit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "exit app"),
		),
	}
}

type keyMap struct {
	page Page

	Save         key.Binding
	ToggleSearch key.Binding
	Quit         key.Binding
	Exit         key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	switch k.page {
	case PageWrite:
		return []key.Binding{k.Save, k.Exit}
	case PageView:
		return []key.Binding{k.Quit, k.Exit}
	case PageSearch:
		return []key.Binding{k.ToggleSearch, k.Exit}
	}
	return nil
}

func (k keyMap) FullHelp() [][]key.Binding {
	return nil
}
