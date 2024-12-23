package models

import (
	"github.com/Gabryjiel/git_config_manager/git"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ValueScope int

const (
	SCOPE_LOCAL ValueScope = iota
	SCOPE_GLOBAL
	SCOPE_SYSTEM
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
	case git.GitConfigProp:
		this.prop = msg
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return this, ExitProgram()
		case tea.KeyEsc:
			return this, SwitchModel(MODEL_SEARCH)
		case tea.KeyRunes:
			this.textField.Value += msg.String()
		case tea.KeyBackspace:
			this.textField.removeLastCharacter()
		case tea.KeyCtrlW:
			this.textField.removeLastWord()
		case tea.KeyDown:
			fallthrough
		case tea.KeyCtrlN:
			this.cursor.moveDown(2)
		case tea.KeyUp:
			fallthrough
		case tea.KeyCtrlP:
			this.cursor.moveUp()
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

	output += renderPropValue(SCOPE_LOCAL, this.prop.Values.Local, this.cursor)
	output += renderPropValue(SCOPE_GLOBAL, this.prop.Values.Global, this.cursor)
	output += renderPropValue(SCOPE_SYSTEM, this.prop.Values.System, this.cursor)

	return output
}

func renderPropValue(scope ValueScope, value string, cursor MenuCursor) string {
	defaultStyle := lipgloss.NewStyle()
	emptyStyle := lipgloss.NewStyle().Width(40)
	valueStyle := lipgloss.NewStyle().Width(40)
	output := ""
	label := "Default label"

	switch scope {
	case SCOPE_LOCAL:
		label = "Local"
		valueStyle = emptyStyle.Foreground(localValueStyle.GetForeground())
	case SCOPE_GLOBAL:
		label = "Global"
		valueStyle = emptyStyle.Foreground(globalValueStyle.GetForeground())
	case SCOPE_SYSTEM:
		label = "System"
		valueStyle = emptyStyle.Foreground(systemValueStyle.GetForeground())
	}

	if value == "" {
		value = "..."
	}

	if int(scope) == int(cursor) {
		defaultStyle = defaultStyle.Background(lipgloss.ANSIColor(8))
	}

	left := emptyStyle.Align(lipgloss.Left).Render(label)
	right := valueStyle.Align(lipgloss.Right).Render(value)

	output += defaultStyle.Padding(0, 1).Render(left+right) + "\n"

	return output
}
