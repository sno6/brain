package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type indexModel struct {
	actions list.Model
}

func newIndexModel() *indexModel {
	listItems := []list.Item{
		actionItem{
			title:       "Write",
			description: "Create a new cell",
			page:        PageWrite,
		},
		actionItem{
			title:       "Read",
			description: "Search and view contents of a cell",
			page:        PageSearch,
		},
		actionItem{
			title:       "Learn",
			description: "Spaced repetition learning",
		},
		actionItem{
			title:       "Help",
			description: "View all available commands",
		},
	}

	actions := list.New(listItems, actionDelegate{}, 60, len(listItems)+3)
	actions.Title = "Brain 🧠"
	actions.Styles.Title = titleStyle
	actions.Styles.TitleBar = lipgloss.NewStyle().MarginBottom(1)
	actions.SetShowTitle(true)
	actions.SetShowPagination(false)
	actions.SetFilteringEnabled(false)
	actions.SetShowStatusBar(false)
	actions.SetShowHelp(false)
	actions.KeyMap.NextPage = key.NewBinding()
	actions.KeyMap.PrevPage = key.NewBinding()

	return &indexModel{actions: actions}
}

func (s *indexModel) Init() tea.Cmd {
	return nil
}

func (s *indexModel) View() string {
	return s.actions.View()
}

func (s *indexModel) Update(msg tea.Msg) (*indexModel, tea.Cmd) {
	var cmd tea.Cmd
	s.actions, cmd = s.actions.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			a, ok := s.actions.SelectedItem().(actionItem)
			if !ok {
				return s, cmd
			}

			cmd = tea.Batch(cmd, changePage(a.page))
		}
	}

	return s, cmd
}

type actionItem struct {
	title, description string
	page               Page
}

func (actionItem) Description() string { return "" }
func (actionItem) FilterValue() string { return "" }

type actionDelegate struct{}

func (d actionDelegate) Height() int                             { return 1 }
func (d actionDelegate) Spacing() int                            { return 0 }
func (d actionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d actionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(actionItem)
	if !ok {
		return
	}

	data := lipgloss.NewStyle().Bold(true).Render(item.title)
	if item.description != "" {
		data = data + " • " + item.description
	}

	var cursor string
	if index == m.Index() {
		cursor = cursorStyle.Render("➜ ")
		data = selectedItemStyle.Render(data)
	} else {
		data = "  " + data
	}

	fmt.Fprintf(w, "%s%s\n", cursor, data)
}
