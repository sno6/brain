package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sno6/brain"
)

// Page is an enum type that we use to determine what
// components to render. This type can be sent as a tea.Msg
// from sub-components to change the page.
type Page uint8

const (
	PageSearch Page = iota
	PageWrite
	PageView
)

func changePage(p Page) func() tea.Msg {
	return func() tea.Msg {
		return p
	}
}

// Base app styling for the whole user interface.
var appStyle = lipgloss.
	NewStyle().
	Padding(1, 2)

// App is the entrypoint for the Brain UI.
type App struct {
	brain *brain.Brain

	curPage  Page
	search   *searchModel
	cellView *cellViewModel
	cellList *cellListModel
}

// NewApp returns a new tea.Model with all sub models.
func NewApp(brain *brain.Brain, startingPage Page) *App {
	return &App{
		brain:    brain,
		curPage:  startingPage,
		search:   newSearchModel(),
		cellView: newCellViewModel(startingPage == PageWrite),
		cellList: newCellListModel(),
	}
}

// Init initialises all sub models.
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		a.search.Init(),
		a.cellList.Init(),
		a.cellView.Init(),
	)
}

// View renders the app by rendering all sub models.
func (a *App) View() string {
	if a.curPage == PageSearch {
		return appStyle.Render(
			lipgloss.JoinVertical(0, a.cellList.View(), a.search.View()),
		)
	}

	return appStyle.Render(a.cellView.View())
}

// Start runs the Brain user interface.
func (a *App) Start() error {
	return tea.NewProgram(a).Start()
}

// Update is the main app update loop, it updates all sub models,
// and handles any and all calls out to the Brain interface.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Exit out of the UI on ctrl+x and esc key presses.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return a, tea.Quit
		}
	}

	// A sub-component has triggered a page change.
	if p, ok := msg.(Page); ok {
		a.curPage = p
	}

	cmd := a.updateSubModels(msg)

	// The user has just stopped typing a query in the search bar.
	// Send what they typed to Brain and create a tea.Cmd for the results,
	// so that cells can listen and display the findings.
	if s, ok := msg.(searchMessage); ok {
		cmd = tea.Batch(cmd, a.searchBrain(string(s)))
	}

	return a, cmd
}

func (a *App) updateSubModels(msg tea.Msg) tea.Cmd {
	var (
		queryCmd, cellListCmd, cellViewCmd tea.Cmd
	)

	switch a.curPage {
	case PageSearch:
		a.search, queryCmd = a.search.Update(msg)
		a.cellList, cellListCmd = a.cellList.Update(msg)
	case PageView, PageWrite:
		a.cellView, cellViewCmd = a.cellView.Update(msg)
	}

	return tea.Batch(queryCmd, cellListCmd, cellViewCmd)
}

type listItems []*brain.Cell

func (a *App) searchBrain(search string) func() tea.Msg {
	return func() tea.Msg {
		cells, _ := a.brain.List(search)
		return listItems(cells)
	}
}
