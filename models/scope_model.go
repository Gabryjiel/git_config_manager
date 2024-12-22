package models

import (
	"github.com/Gabryjiel/git_config_manager/git"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ScopeModel struct {
	tea.Model
	prop      git.GitConfigProp
	textField TextField
	values    git.GitConfigEntryValues
	cursor    MenuCursor
}

func NewScopeModel() ScopeModel {
	return ScopeModel{}
}

func (this ScopeModel) Init() tea.Cmd {
	return nil
}

func (this ScopeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return this, tea.Quit
		case tea.KeyEsc:
			return NewSearchModel(), nil
		}
	}
	return this, nil
}

func (this ScopeModel) View() string {
	output := ""

	output += renderHeader()
	output += arrowCinStyle.Render(" Value: ")
	output += this.textField.Value
	output += lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15)).Render("█")
	output += "\n\n"

	output += "Local :"
	output += "Global:"
	output += "System:"

	return output
}
