package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sno6/brain"
)

// Base app styling for the whole user interface.
var appStyle = lipgloss.
	NewStyle().
	Padding(1, 2)

// App is the entrypoint for the Brain UI.
type App struct {
	brain *brain.Brain

	search *searchModel
	cells  *cellListModel
}

// NewApp returns a new tea.Model with all sub models.
func NewApp(brain *brain.Brain) *App {
	return &App{
		brain:  brain,
		search: newSearchModel(),
		cells:  newCellListModel(),
	}
}

// Init initialises all sub models.
func (a *App) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, a.search.Init(), a.cells.Init())
}

// View renders the app by rendering all sub models.
func (a *App) View() string {
	return appStyle.Render(
		a.cells.View() + "\n" + a.search.View(),
	)
}

// Start runs the Brain user interface.
func (a *App) Start() error {
	return tea.NewProgram(a).Start()
}

// Update is the main app update loop, it updates all sub models,
// and handles any and all calls out to the Brain interface.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := a.updateSubModels(msg)

	// Exit out of the UI on ctrl+x and esc key presses.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return a, tea.Quit
		}
	}

	// The user has just stopped typing a query in the search bar.
	// Send what they typed to Brain and create a tea.Cmd for the results,
	// so that cells can listen and display the findings.
	if s, ok := msg.(searchMessage); ok {
		cmd = tea.Batch(cmd, a.searchBrain(string(s)))
	}

	return a, cmd
}

func (a *App) updateSubModels(msg tea.Msg) tea.Cmd {
	var queryCmd, cellsCmd tea.Cmd
	a.search, queryCmd = a.search.Update(msg)
	a.cells, cellsCmd = a.cells.Update(msg)
	return tea.Batch(queryCmd, cellsCmd)
}

type listItems []*brain.Cell

func (a *App) searchBrain(search string) func() tea.Msg {
	return func() tea.Msg {
		cells, _ := a.brain.List(search)
		return listItems(cells)
	}
}
