package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sno6/brain"
)

// Base app styling for the whole user interface.
var appStyle = lipgloss.NewStyle().Padding(1, 2)

// App is the entrypoint for the Brain UI.
type App struct {
	brain *brain.Brain

	curPage                     Page
	width, height               int
	index                       *indexModel
	search                      *searchModel
	cellList                    *cellListModel
	writeCellView, readCellView *cellViewModel
}

// NewApp returns a new tea.Model with all sub models.
func NewApp(brain *brain.Brain, startingPage Page) *App {
	return &App{
		brain:         brain,
		curPage:       startingPage,
		index:         newIndexModel(),
		search:        newSearchModel(),
		cellList:      newCellListModel(),
		writeCellView: newCellViewModel(PageWrite),
		readCellView:  newCellViewModel(PageView),
	}
}

// Init initialises all sub models.
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		a.index.Init(),
		a.search.Init(),
		a.cellList.Init(),
		a.readCellView.Init(),
		a.writeCellView.Init(),
	)
}

// View renders the app by rendering all sub models.
func (a *App) View() string {
	switch a.curPage {
	case PageIndex:
		return appStyle.Render(a.index.View())
	case PageSearch:
		s := lipgloss.JoinVertical(0, a.cellList.View(), a.search.View())
		return appStyle.Render(s)
	case PageWrite:
		return appStyle.Render(a.writeCellView.View())
	case PageView:
		return appStyle.Render(a.readCellView.View())
	}
	return "<unknown page>"
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
	case tea.WindowSizeMsg:
		a.propagateDimensions(msg.Width, msg.Height)
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

	// The user has clicked ctrl+s on the write cell page.
	if c, ok := msg.(savedCell); ok {
		a.brain.Write(string(c))
		return a, tea.Quit
	}

	cmd := a.updateSubModels(msg)

	// The user has just stopped typing a query in the search bar.
	// Send what they typed to Brain and create a tea.Cmd for the results,
	// so that cells can listen and display the findings.
	if s, ok := msg.(searchMessage); ok {
		cmd = tea.Batch(cmd, a.searchBrain(s))
	}

	return a, cmd
}

func (a *App) updateSubModels(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch a.curPage {
	case PageIndex:
		a.index, cmd = a.index.Update(msg)
	case PageWrite:
		a.writeCellView, cmd = a.writeCellView.Update(msg)
	case PageView:
		a.readCellView, cmd = a.readCellView.Update(msg)
	case PageSearch:
		var searchCmd, cellListCmd tea.Cmd
		a.search, searchCmd = a.search.Update(msg)
		a.cellList, cellListCmd = a.cellList.Update(msg)
		cmd = tea.Batch(cmd, searchCmd, cellListCmd)
	}

	return cmd
}

func (a *App) propagateDimensions(width, height int) {
	if a.width == width && a.height == height {
		return
	}

	a.search.setDimensions(width, height)
	a.writeCellView.setDimensions(width, height)
	a.readCellView.setDimensions(width, height)
	a.cellList.setDimensions(width, height)
}

type listItems []*brain.Cell

func (a *App) searchBrain(sm searchMessage) func() tea.Msg {
	return func() tea.Msg {
		cells, _ := a.brain.List(sm.val, sm.mode)
		return listItems(cells)
	}
}
