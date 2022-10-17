package tui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	textCursorStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("212"))

	focusedCursorLineStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("57")).
				Foreground(lipgloss.Color("230"))

	focusedStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238"))
)

type cellViewModel struct {
	text textarea.Model
	help *helpModel
	page Page
}

func newCellViewModel(page Page) *cellViewModel {
	text := textarea.New()

	text.Prompt = ""
	text.Cursor.Style = textCursorStyle
	text.FocusedStyle.CursorLine = focusedCursorLineStyle
	text.FocusedStyle.Base = focusedStyle
	text.BlurredStyle.Base = text.FocusedStyle.Base
	text.CharLimit = -1

	if page == PageWrite {
		text.Focus()
	} else {
		text.Blur()
	}

	return &cellViewModel{
		text: text,
		help: newHelpModel(page),
		page: page,
	}
}

func (c *cellViewModel) Init() tea.Cmd {
	return c.help.Init()
}

// View renders the app by rendering all sub models.
func (c *cellViewModel) View() string {
	title := titleStyle.Render("Brain 🧠")
	return lipgloss.JoinVertical(0, title, c.text.View(), c.help.View())
}

func (c *cellViewModel) Update(msg tea.Msg) (*cellViewModel, tea.Cmd) {
	// Exit out of the full cell view on 'q' keypress.
	if c.page == PageView {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyRunes:
				if msg.String() == "q" {
					// TODO: prev page?
					return c, changePage(PageSearch)
				}
			}
		}
	}

	if c.page == PageWrite {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlS: // CMD?
				return c, saveCell(savedCell(c.text.Value()))
			}
		}
	}

	if s, ok := msg.(viewCellMessage); ok {
		c.text.SetValue(string(s[11:]))
	}

	var helpCmd, textCmd tea.Cmd
	c.help, helpCmd = c.help.Update(msg)
	c.text, textCmd = c.text.Update(msg)
	return c, tea.Batch(helpCmd, textCmd)
}

func (c *cellViewModel) setDimensions(width, height int) {
	c.text.SetWidth(width - 5)
	c.text.SetHeight(int(float64(height) * 0.7))
}
