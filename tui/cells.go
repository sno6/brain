package tui

import (
	"fmt"
	"io"
	"strings"
	"time"

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
}

func newCellListModel() *cellListModel {
	cells := list.New(nil, cellDelegate{}, 60, 0)
	cells.Title = "Cells"
	cells.Styles.TitleBar = titleBarStyle
	cells.Styles.Title = titleStyle
	cells.Paginator.PerPage = 10
	cells.Styles.PaginationStyle.PaddingBottom(1)
	cells.DisableQuitKeybindings()
	cells.SetShowTitle(true)
	cells.SetFilteringEnabled(false)
	cells.SetShowStatusBar(false)
	cells.SetShowHelp(false)

	return &cellListModel{cells: cells}
}

func (c *cellListModel) Init() tea.Cmd {
	return nil
}

func (c *cellListModel) Update(msg tea.Msg) (*cellListModel, tea.Cmd) {
	cmd := c.updateSubModels(msg)

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
	cells := make([]list.Item, len(items))
	for i, item := range items {
		cells[i] = cell{
			data: string(item.Data()),
			ts:   item.Timestamp(),
		}
	}

	padding := 4
	c.cells.SetItems(cells)
	c.cells.SetHeight(max(len(cells), 10) + padding)
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
