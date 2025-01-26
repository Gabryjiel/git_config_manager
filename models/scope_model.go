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
	cursor    MenuCursor
	isEditing bool
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
			return this, tea.Quit
		case tea.KeyEsc:
			if this.isEditing {
				this.isEditing = false
			} else {
				return NewSearchModel(), nil
			}
		case tea.KeyRunes:
			if this.isEditing {
				this.textField.Value += msg.String()
			}
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
		case tea.KeyEnter:
			this.isEditing = !this.isEditing

			if !this.isEditing {
			}
		}
	}
	return this, nil
}

func (this ScopeModel) View() string {
	output := ""

	output += renderHeader()
	output += arrowCinStyle.Render(" Value: ")
	output += this.textField.Value

	if this.isEditing {
		output += lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15)).Render("█")
	}

	output += "\n\n"

	localValue, _ := this.prop.Values["local"]
	globalValue, _ := this.prop.Values["global"]
	systemValue, _ := this.prop.Values["system"]

	output += renderPropValue(git.SCOPE_LOCAL, localValue, this.cursor)
	output += renderPropValue(git.SCOPE_GLOBAL, globalValue, this.cursor)
	output += renderPropValue(git.SCOPE_SYSTEM, systemValue, this.cursor)

	return output
}

func renderPropValue(scope git.ValueScope, value string, cursor MenuCursor) string {
	defaultStyle := lipgloss.NewStyle()
	emptyStyle := lipgloss.NewStyle().Width(40)
	valueStyle := lipgloss.NewStyle().Width(40)
	output := ""
	label := "Default label"

	switch scope {
	case git.SCOPE_LOCAL:
		label = "Local"
		valueStyle = emptyStyle.Foreground(localValueStyle.GetForeground())
	case git.SCOPE_GLOBAL:
		label = "Global"
		valueStyle = emptyStyle.Foreground(globalValueStyle.GetForeground())
	case git.SCOPE_SYSTEM:
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
