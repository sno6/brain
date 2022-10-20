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

	deleteQuestionStyle = lipgloss.NewStyle().MarginLeft(1)
)

type cellViewModel struct {
	text     textarea.Model
	help     *helpModel
	editable bool

	// The ID of the document in our index that we are
	// currently referencing.
	currentDocID string

	// The user has clicked 'x' on a cell, we track this state
	// so we can hide the help and present a "Are you sure?" message.
	deleteDialouge bool
	deleteOption   bool
}

func (c *cellViewModel) setEditable(e bool) {
	c.editable = e

	if e {
		c.text.Focus()
		c.help = newHelpModel(PageWrite)
	} else {
		c.text.Blur()
		c.help = newHelpModel(PageView)
	}
}

func newCellViewModel() *cellViewModel {
	text := textarea.New()

	text.Prompt = ""
	text.Cursor.Style = textCursorStyle
	text.FocusedStyle.CursorLine = focusedCursorLineStyle
	text.FocusedStyle.Base = focusedStyle
	text.BlurredStyle.Base = text.FocusedStyle.Base
	text.CharLimit = -1
	text.Blur()

	return &cellViewModel{
		text: text,
		help: newHelpModel(PageView),
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

	if c.deleteDialouge {
		var options string
		if c.deleteOption {
			options = textCursorStyle.Bold(true).Render("yes") + " / no"
		} else {
			options = "yes / " + textCursorStyle.Bold(true).Render("no")

		}

		views = append(views,
			deleteQuestionStyle.Render(fmt.Sprintf("Are you sure? %s", options)))
	} else {
		views = append(views, c.help.View())
	}

	return lipgloss.JoinVertical(0, views...)
}

func (c *cellViewModel) Update(msg tea.Msg) (*cellViewModel, tea.Cmd) {
	if !c.editable {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyTab, tea.KeyRight, tea.KeyLeft:
				if !c.deleteDialouge {
					break
				}
				c.deleteOption = !c.deleteOption
			case tea.KeyEnter:
				if !c.deleteDialouge {
					break
				}

				if shouldDelete := c.deleteOption; shouldDelete {
					docID := c.currentDocID
					c.currentDocID = ""
					c.deleteDialouge = false

					return c, deleteCell(docID)
				} else {
					c.deleteDialouge = false
					c.deleteOption = true
				}
			case tea.KeyRunes:
				msgStr := msg.String()

				if msgStr == "q" {
					c.currentDocID = ""
					c.deleteDialouge = false
					return c, changePage(PageSearch)
				}
				if msgStr == "x" {
					c.deleteDialouge = true
					c.deleteOption = true
				}
				if msgStr == "e" {
					c.setEditable(true)
					return c, changePage(PageWrite)
				}
			}
		}
	}

	if c.editable {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlS:
				docID := c.currentDocID
				c.currentDocID = ""

				return c, saveCell(docID, c.text.Value())
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

func (c *cellViewModel) setDimensions(width, height int) {
	c.text.SetWidth(width - 5)
	c.text.SetHeight(int(float64(height) * 0.7))
}
