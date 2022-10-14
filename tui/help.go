package tui

import tea "github.com/charmbracelet/bubbletea"

type helpModel struct{}

func newHelpModel() *helpModel {
	return &helpModel{}
}

func (h *helpModel) Init() tea.Cmd {
	return nil
}

func (h *helpModel) Update(msg tea.Msg) (*helpModel, tea.Cmd) {
	return h, nil
}

func (h *helpModel) View() string {
	return "<help>ajfdajksdf"
}
