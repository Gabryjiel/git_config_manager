package models

import (
	"github.com/Gabryjiel/git_config_manager/git"
	"github.com/Gabryjiel/git_config_manager/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ScopeModel struct {
	prop      git.GitConfigProp
	input     TextInputModel
	cursor    MenuCursor
	isEditing bool
}

func NewScopeModel() *ScopeModel {
	return &ScopeModel{
		input:     NewTextField(),
		cursor:    0,
		isEditing: false,
	}
}

func (this *ScopeModel) Init() tea.Cmd {
	return nil
}

func (this *ScopeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case GitConfigSetResult:
		this.isEditing = false
		this.prop.Values[this.getScopeFromCursor()] = this.input.GetValue()
		this.input.Clear()
		return this, nil
	case git.GitConfigProp:
		this.prop = msg
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			this.isEditing = false
			return this, SwitchSubmodel(MODEL_ID_SEARCH)
		case tea.KeyDown:
			this.cursor.moveDown(2)
		case tea.KeyUp:
			this.cursor.moveUp()
		case tea.KeyEnter:
			if this.isEditing {
				this.isEditing = false
				return this, CmdGitConfigSet(this.prop.GetName(), this.input.GetValue())
			} else {
				this.isEditing = true
				value, ok := this.prop.Values[this.getScopeFromCursor()]
				if ok {
					this.input.SetValue(value)
				}
				return this, nil
			}
		}
	}

	if this.isEditing {
		_, cmd := this.input.Update(msg)
		return this, cmd
	}
	return this, nil
}

func (this *ScopeModel) View() string {
	output := ""

	output += renderHeader()
	output += arrowCinStyle.Render(" Value: ")
	output += this.input.View()
	output += "\n\n"

	localValue, _ := this.prop.Values["local"]
	globalValue, _ := this.prop.Values["global"]
	systemValue, _ := this.prop.Values["system"]

	output += renderPropValue(git.SCOPE_LOCAL, localValue, this.cursor)
	output += renderPropValue(git.SCOPE_GLOBAL, globalValue, this.cursor)
	output += renderPropValue(git.SCOPE_SYSTEM, systemValue, this.cursor)

	return output
}

// Helpers

func (this *ScopeModel) getScopeFromCursor() string {
	switch this.cursor {
	case 0:
		return "local"
	case 1:
		return "global"
	case 2:
		return "system"
	default:
		return "local"
	}
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

// Commands

func CmdScopeSelectProp(prop git.GitConfigProp) tea.Cmd {
	return func() tea.Msg {
		return prop
	}
}

type GitConfigSetResult struct {
	result  bool
	message string
}

func CmdGitConfigSet(name string, value string) tea.Cmd {
	return func() tea.Msg {
		content, err := utils.ExecuteCommand("git", "config", "set", name, value)

		if err != nil {
			return GitConfigSetResult{
				message: err.Error(),
				result:  false,
			}
		}

		return GitConfigSetResult{
			result:  true,
			message: content,
		}
	}
}
