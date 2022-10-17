package tui

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			MarginLeft(0).
			Padding(0, 2).
			Italic(true).
			Bold(true).
			Foreground(lipgloss.Color("#FFF")).
			Background(lipgloss.Color("#F25D94"))

	itemDateStyle = lipgloss.
			NewStyle().
			Bold(true)

	selectedItemStyle = lipgloss.
				NewStyle().
				Foreground(lipgloss.Color("#FFF"))

	cursorStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"})
)

// A cellListModel is a wrapper around a list.Model that displays cells.
type cellListModel struct {
	cells list.Model
	page  int
}

func newCellListModel() *cellListModel {
	cells := list.New(nil, cellDelegate{}, 60, 0)
	cells.Title = "Brain 🧠"
	cells.Styles.Title = titleStyle
	cells.Styles.TitleBar = lipgloss.NewStyle().MarginBottom(1)
	cells.Paginator.PerPage = 10
	cells.Styles.PaginationStyle.PaddingBottom(1)
	cells.SetShowTitle(true)
	cells.SetFilteringEnabled(false)
	cells.SetShowStatusBar(false)
	cells.SetShowHelp(false)

	// Disable next/prev pagination.
	cells.KeyMap.NextPage = key.NewBinding()
	cells.KeyMap.PrevPage = key.NewBinding()
	cells.DisableQuitKeybindings()

	return &cellListModel{cells: cells}
}

func (c *cellListModel) Init() tea.Cmd {
	return nil
}

func (c *cellListModel) Update(msg tea.Msg) (*cellListModel, tea.Cmd) {
	cmd := c.updateSubModels(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			c, ok := c.cells.SelectedItem().(cell)
			if ok {
				cmd = tea.Batch(cmd, viewCellCommand(c.data), changePage(PageView))
			}
		}
	}

	// Figure out if the page has changed so we can update the list size.
	if c.page != c.cells.Paginator.Page {
		c.page = c.cells.Paginator.Page

		padding := 3
		itemLen := len(c.cells.Items())
		pageLen := c.cells.Paginator.ItemsOnPage(itemLen)
		c.cells.SetHeight(pageLen + padding)
	}

	// Search has found some new items, we need to update
	// our internal model and render the list items.
	if items, ok := msg.(listItems); ok {
		c.updateListItems(items)
	}

	return c, cmd
}

func (c *cellListModel) updateSubModels(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	c.cells, cmd = c.cells.Update(msg)
	return cmd
}

func (c *cellListModel) updateListItems(items listItems) {
	if len(items) == 0 {
		c.cells.SetItems(nil)
		c.cells.SetHeight(0)
		return
	}

	cells := make([]list.Item, len(items))
	for i, item := range items {
		cells[i] = cell{
			data: item.Data(),
			ts:   item.Timestamp(),
		}
	}

	padding := 3
	c.cells.SetItems(cells)
	c.cells.SetHeight(len(cells) + padding)
}

func (c *cellListModel) View() string {
	return c.cells.View()
}

// A cell is the UI element for a row in the list.
type cell struct {
	data string
	ts   time.Time
}

func (c cell) Description() string { return c.data }
func (c cell) FilterValue() string { return "" }

// A cellDelegate is responsible for updating and rendering an individual
// cell in a list.Model.
type cellDelegate struct{}

func (d cellDelegate) Height() int                             { return 1 }
func (d cellDelegate) Spacing() int                            { return 0 }
func (d cellDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d cellDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(cell)
	if !ok {
		return
	}

	date := itemDateStyle.Render(fmt.Sprintf("%02d/%02d/%02d", item.ts.Day(), item.ts.Month(), item.ts.Year()))
	data := preview(item.data)

	var cursor string
	if index == m.Index() {
		cursor = cursorStyle.Render("➜ ")
		data = "  " + selectedItemStyle.Render(data)
	} else {
		date = "  " + date
		data = "  " + data
	}

	fmt.Fprintf(w, "%s%s\n%s", cursor, date, data)
}

// How many characters to show for long cells.
const previewLength = 70

func preview(data string) string {
	// Exclude the date portion of the cell.
	data = data[11:]

	if len(data) < previewLength {
		return data
	}

	var prev string
	if nl := strings.Index(data, "\n"); nl > -1 {
		if nl > previewLength {
			prev = data[:previewLength]
		} else {
			prev = data[:nl]
		}
	}

	return prev + "..."
}
