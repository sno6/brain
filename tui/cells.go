package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleBarStyle = lipgloss.
			NewStyle().
			Padding(0, 0, 1)

	titleStyle = lipgloss.
			NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			Padding(0, 1)
)

type cellListModel struct {
	cells list.Model
}

func newCellListModel() *cellListModel {
	return &cellListModel{
		cells: initList(),
	}
}

func (c *cellListModel) Init() tea.Cmd {
	return nil
}

func (c *cellListModel) Update(msg tea.Msg) (*cellListModel, tea.Cmd) {
	// We have found some new list items, render them.
	if items, ok := msg.(listItems); ok {
		c.updateList(items)
	}

	cmd := c.updateSubModels(msg)
	return c, cmd
}

func (c *cellListModel) updateSubModels(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	c.cells, cmd = c.cells.Update(msg)
	return cmd
}

func (c *cellListModel) updateList(items listItems) {
	cells := make([]list.Item, len(items))
	for i, item := range items {
		cells[i] = cell{
			title:       fmt.Sprintf("Element %d", i+1),
			description: item,
		}
	}

	c.cells.SetItems(cells)
}

func (c *cellListModel) View() string {
	return c.cells.View()
}

// A cell is the UI element for a row in the list.
type cell struct {
	title       string
	description string
}

func (c cell) Title() string       { return c.title }
func (c cell) Description() string { return c.description }
func (c cell) FilterValue() string { return c.title }

func initList() list.Model {
	itemDelegate := list.NewDefaultDelegate()
	cells := list.New(nil, itemDelegate, 100, 0)

	cells.Title = "Cells"
	cells.Styles.TitleBar = titleBarStyle
	cells.Styles.Title = titleStyle

	cells.SetShowTitle(true)
	cells.SetFilteringEnabled(false)
	cells.SetShowStatusBar(false)
	cells.SetShowHelp(false)

	return cells
}
