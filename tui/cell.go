package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cellViewModel struct {
	text textarea.Model
	help *helpModel

	editable bool
}

func newCellViewModel(editable bool) *cellViewModel {
	text := textarea.New()
	text.SetWidth(156)
	text.SetHeight(35)

	text.Prompt = ""
	text.Cursor.Style = lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("212"))

	text.FocusedStyle.CursorLine = lipgloss.NewStyle().
		Background(lipgloss.Color("57")).
		Foreground(lipgloss.Color("230"))

	text.FocusedStyle.Base = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("238"))

	text.BlurredStyle.Base = text.FocusedStyle.Base

	if editable {
		text.Focus()
	} else {
		text.Blur()
	}

	return &cellViewModel{
		text:     text,
		help:     newHelpModel(),
		editable: editable,
	}
}

func (c *cellViewModel) Init() tea.Cmd {
	return c.help.Init()
}

// View renders the app by rendering all sub models.
func (c *cellViewModel) View() string {
	var doc strings.Builder

	doc.WriteString("🧠 Brain\n")
	doc.WriteString(c.text.View())
	doc.WriteString(c.help.View())

	return doc.String()
}

func (c *cellViewModel) Update(msg tea.Msg) (*cellViewModel, tea.Cmd) {
	// Exit out of the full cell view on 'q' keypress.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			if msg.String() == "q" {
				return c, changePage(PageSearch)
			}
		}
	}

	if s, ok := msg.(viewCellMessage); ok {
		c.text.SetValue(string(s))
	}

	var helpCmd, textCmd tea.Cmd
	c.help, helpCmd = c.help.Update(msg)
	c.text, textCmd = c.text.Update(msg)
	return c, tea.Batch(helpCmd, textCmd)
}
