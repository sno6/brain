package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	textCursorStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("212"))

	focusedCursorLineStyle = lipgloss.
				NewStyle().
				Background(lipgloss.Color("#684EFF")).
				Foreground(lipgloss.Color("230"))

	focusedStyle = lipgloss.
			NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238"))

	deleteQuestionStyle = lipgloss.
				NewStyle().
				MarginLeft(1)

	questionSelectedStyle = textCursorStyle.Bold(true)
)

type cellViewModel struct {
	text     textarea.Model
	help     *helpModel
	editable bool

	// The ID of the document that we are currently viewing.
	currentDocID string

	// The user has clicked 'x' on a cell, we track this state
	// so we can hide the help and present an "Are you sure?" message.
	deleteDialogOpen bool
	deleteOption     bool
}

func newCellViewModel() *cellViewModel {
	text := textarea.New()
	text.Prompt = ""
	text.Cursor.Style = textCursorStyle
	text.FocusedStyle.CursorLine = focusedCursorLineStyle
	text.FocusedStyle.Base = focusedStyle
	text.BlurredStyle.Base = text.FocusedStyle.Base
	text.CharLimit = -1

	return &cellViewModel{
		text:         text,
		help:         newHelpModel(PageView),
		deleteOption: true,
	}
}

func (c *cellViewModel) Init() tea.Cmd {
	return c.help.Init()
}

// View renders the app by rendering all sub models.
func (c *cellViewModel) View() string {
	views := []string{
		titleStyle.Render("Brain 🧠"),
		c.text.View(),
	}

	if c.deleteDialogOpen {
		views = append(views, c.renderDeleteDialog())
	} else {
		views = append(views, c.help.View())
	}

	return lipgloss.JoinVertical(0, views...)
}

func (c *cellViewModel) Update(msg tea.Msg) (*cellViewModel, tea.Cmd) {
	if c.editable {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlS:
				id := c.currentDocID
				c.reset()

				return c, saveCell(id, c.text.Value())
			}
		}
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyTab, tea.KeyRight, tea.KeyLeft:
				if !c.deleteDialogOpen {
					break
				}

				c.deleteOption = !c.deleteOption
			case tea.KeyEnter:
				if !c.deleteDialogOpen {
					break
				}

				if c.deleteOption {
					id := c.currentDocID
					c.reset()

					return c, deleteCell(id)
				} else {
					// User selected "no", close dialog and reset option to "yes".
					c.deleteDialogOpen = false
					c.deleteOption = true
				}
			case tea.KeyRunes:
				switch msg.String() {
				case "x":
					c.deleteDialogOpen = true
				case "e":
					c.setEditable(true)
					return c, changePage(PageWrite)
				case "q":
					c.reset()
					return c, changePage(PageSearch)
				}
			}
		}
	}

	if s, ok := msg.(viewCellMessage); ok {
		c.currentDocID = s.id
		c.text.SetValue(s.content)
	}

	var helpCmd, textCmd tea.Cmd
	c.help, helpCmd = c.help.Update(msg)
	c.text, textCmd = c.text.Update(msg)
	return c, tea.Batch(helpCmd, textCmd)
}

func (c *cellViewModel) renderDeleteDialog() string {
	var options string
	if c.deleteOption {
		options = questionSelectedStyle.Render("yes") + " / no"
	} else {
		options = "yes / " + questionSelectedStyle.Render("no")
	}

	return deleteQuestionStyle.Render(fmt.Sprintf("Are you sure? %s", options))
}

func (c *cellViewModel) setDimensions(width, height int) {
	c.text.SetWidth(width - 5)
	c.text.SetHeight(int(float64(height) * 0.7))
}

func (c *cellViewModel) setEditable(e bool) {
	if e {
		c.text.Focus()
		c.help.setPage(PageWrite)
		c.editable = true
	} else {
		c.text.Blur()
		c.help.setPage(PageView)
		c.editable = false
	}
}

func (c *cellViewModel) reset() {
	c.currentDocID = ""
	c.deleteDialogOpen = false
	c.deleteOption = true
}
