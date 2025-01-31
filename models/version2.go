package models

import (
	"slices"

	"github.com/Gabryjiel/git_config_manager/git"
	"github.com/Gabryjiel/git_config_manager/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GitScope int

const (
	GIT_SCOPE_LOCAL GitScope = iota
	GIT_SCOPE_GLOBAL
	GIT_SCOPE_SYSTEM
)

type MainModel struct {
	isExiting     bool
	isEditing     bool
	searchInput   TextInputModel
	valueInput    TextInputModel
	props         []git.GitConfigProp
	filteredProps []git.GitConfigProp
	scope         GitScope
	cursor        int
	message       string
	listStart     int
	onlyWithValue bool
}

func (this *MainModel) Init() tea.Cmd {
	return Cmd_GetGitConfigEntries
}

func (this *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case Msg_GitConfigSetResult:
		if !msg.result {
			return this, Cmd_DisplayMessage(msg.message)
		}

		this.isEditing = false
		index := slices.IndexFunc(this.props, func(prop git.GitConfigProp) bool {
			return prop.GetName() == msg.name
		})

		this.props[index].Values[msg.scope] = msg.value
		this.filteredProps[this.cursor].Values[msg.scope] = msg.value
		this.valueInput.Clear()

		return this, Cmd_DisplayMessage("Changed value of " + msg.name + " to " + msg.value)

	case Msg_DisplayMessage:
		this.message = msg.message
		return this, nil

	case []git.GitConfigProp:
		this.props = msg
		this.filteredProps = msg
		return this, nil

	case MsgInputChanged:
		this.filteredProps = git.FilterGitConfigProps(this.props, this.searchInput.GetValue(), this.onlyWithValue)

		if len(this.filteredProps) == 0 {
			this.cursor = 0
		} else if this.cursor >= len(this.filteredProps) {
			this.cursor = len(this.filteredProps) - 1
		}

		return this, nil

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyTab:
			this.onlyWithValue = !this.onlyWithValue
			return this, CmdInputChanged

		case tea.KeyCtrlC:
			this.isExiting = true
			return this, tea.Quit

		case tea.KeyCtrlN:
			fallthrough
		case tea.KeyDown:
			if this.cursor < len(this.filteredProps)-1 {
				this.cursor++
			}
			if this.cursor >= this.listStart+10 {
				this.listStart = this.cursor - 10
			}

		case tea.KeyCtrlP:
			fallthrough
		case tea.KeyUp:
			if this.cursor > 0 {
				this.cursor--
			}
			if this.cursor < this.listStart {
				this.listStart = this.cursor
			}

		case tea.KeyPgDown:
			this.cursor += 10
			this.listStart += 10

			if this.cursor > len(this.props) {
				this.cursor = len(this.props) - 9
			}
			if this.listStart > len(this.props) {
				this.listStart = len(this.props) - 9
			}

		case tea.KeyPgUp:
			this.cursor -= 10
			this.listStart -= 10

			if this.cursor < 0 {
				this.cursor = 0
			}
			if this.listStart < 0 {
				this.listStart = 0
			}

		case tea.KeyCtrlRight:
			if this.scope == 2 {
				this.scope = 0
			} else {
				this.scope++
			}

		case tea.KeyCtrlLeft:
			if this.scope == 0 {
				this.scope = 2
			} else {
				this.scope--
			}

		case tea.KeyEsc:
			if this.isEditing {
				this.isEditing = false
				return this, nil
			} else {
				this.isExiting = true
				return this, tea.Quit
			}

		case tea.KeyEnter:
			if this.isEditing {
				name := this.filteredProps[this.cursor].GetName()
				value := this.valueInput.GetValue()

				return this, Cmd_GitConfigSet(name, value, this.scope)
			}

			this.isEditing = !this.isEditing
			oldValue, ok := this.filteredProps[this.cursor].Values[getScopeFromGitScope(this.scope)]
			if ok {
				this.valueInput.SetValue(oldValue)
			} else {
				this.valueInput.SetValue("")
			}

			return this, nil
		}

	}

	if this.isEditing {
		_, cmd := this.valueInput.Update(msg)
		return this, cmd
	} else {
		_, cmd := this.searchInput.Update(msg)
		return this, cmd
	}
}

func (this *MainModel) View() string {
	output := ""

	if this.isExiting {
		return output
	}

	output += renderEasyHeader(this.scope) + "\n"

	if this.isEditing {
		output += "Value: "
		output += this.valueInput.View() + "\n"

	} else {
		output += "Search: "
		output += this.searchInput.View() + "\n"
	}

	centerStyle := lipgloss.NewStyle().Width(80).AlignHorizontal(lipgloss.Center)

	output += centerStyle.Render(renderGap(80)) + "\n"

	for index, prop := range this.filteredProps {
		if index < this.listStart || index > this.listStart+10 {
			continue
		}

		output += renderProp(prop.GetName(), getValueFromScope(prop, this.scope), getColorFromScope(this.scope), index == this.cursor)
	}

	output += centerStyle.Render(renderGap(80)) + "\n"
	output += "Last message: " + this.message

	return output
}

// Helpers

func CreateNewMainModel() *MainModel {
	return &MainModel{}
}

func renderGap(length int) string {
	result := ""

	for i := 0; i < length; i++ {
		result += "-"
	}

	return result
}

func renderEasyHeader(scope GitScope) string {
	result := ""
	result += lipgloss.NewStyle().
		Width(80).
		AlignHorizontal(lipgloss.Center).
		Render("--- gcm v0.0.1 --- " + git.GetGitVersion() + " --- Scope: " + renderHeaderScope(scope) + " --- ")

	return result
}

func renderHeaderScope(gitScope GitScope) string {
	return lipgloss.NewStyle().Foreground(getColorFromScope(gitScope)).Render(getScopeFromGitScope(gitScope))
}

func renderProp(label, value string, valueColor lipgloss.ANSIColor, isSelected bool) string {
	style := lipgloss.NewStyle()

	if isSelected {
		style = style.Background(lipgloss.ANSIColor(8))
	}

	propLabel := style.PaddingLeft(1).Render(label)
	propValue := style.Width(80 - len(label) - 1).AlignHorizontal(lipgloss.Right).PaddingRight(1).Foreground(valueColor).Render(value)

	return propLabel + propValue + "\n"
}

func getValueFromScope(prop git.GitConfigProp, scope GitScope) string {
	switch scope {
	case GIT_SCOPE_LOCAL:
		return prop.Values["local"]
	case GIT_SCOPE_GLOBAL:
		return prop.Values["global"]
	case GIT_SCOPE_SYSTEM:
		return prop.Values["system"]
	}

	return ""
}

func getColorFromScope(scope GitScope) lipgloss.ANSIColor {
	switch scope {
	case GIT_SCOPE_LOCAL:
		return lipgloss.ANSIColor(10)
	case GIT_SCOPE_GLOBAL:
		return lipgloss.ANSIColor(11)
	case GIT_SCOPE_SYSTEM:
		return lipgloss.ANSIColor(12)
	default:
		return lipgloss.ANSIColor(1)
	}
}

func getScopeFromGitScope(gitScope GitScope) string {
	switch gitScope {
	case GIT_SCOPE_LOCAL:
		return "local"
	case GIT_SCOPE_GLOBAL:
		return "global"
	case GIT_SCOPE_SYSTEM:
		return "system"
	default:
		return "unknown"
	}
}

// Commands

func Cmd_GetGitConfigEntries() tea.Msg {
	return git.GetConfigProps()
}

type Msg_DisplayMessage struct {
	message string
}

func Cmd_DisplayMessage(message string) tea.Cmd {
	return func() tea.Msg {
		return Msg_DisplayMessage{message}
	}
}

type Msg_GitConfigSetResult struct {
	result  bool
	message string
	name    string
	value   string
	scope   string
}

func Cmd_GitConfigSet(name string, value string, gitScope GitScope) tea.Cmd {
	return func() tea.Msg {
		scope := getScopeFromGitScope(gitScope)
		content, err := utils.ExecuteCommand("git", "config", "set", "--"+scope, name, value)

		if err != nil {
			return Msg_GitConfigSetResult{
				message: err.Error(),
				result:  false,
				name:    name,
				value:   value,
				scope:   scope,
			}
		}

		return Msg_GitConfigSetResult{
			result:  true,
			message: content,
			name:    name,
			value:   value,
			scope:   scope,
		}
	}
}
