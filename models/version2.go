package models

import (
	"fmt"
	"slices"

	"github.com/Gabryjiel/git_config_manager/git"
	"github.com/Gabryjiel/git_config_manager/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GitScope string

const (
	GIT_SCOPE_LOCAL  GitScope = "local"
	GIT_SCOPE_GLOBAL GitScope = "global"
	GIT_SCOPE_SYSTEM GitScope = "system"
)

type MainModel struct {
	isExiting     bool
	isEditing     bool
	searchInput   textinput.Model
	valueInput    textinput.Model
	props         []git.GitConfigProp
	filteredProps []git.GitConfigProp
	scope         GitScope
	cursor        int
	message       string
	listStart     int
	onlyWithValue bool
	help          help.Model
}

func (this MainModel) Init() tea.Cmd {
	return Cmd_GetGitConfigEntries()
}

func (this MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case Msg_GitConfigSetResult:
		if !msg.result {
			return this, tea.Batch(
				Cmd_DisplayMessage(msg.message),
				Cmd_AddLog(msg.command),
			)
		}

		this.isEditing = false
		index := slices.IndexFunc(this.props, func(prop git.GitConfigProp) bool {
			return prop.GetName() == msg.name
		})

		this.props[index].Values[msg.scope] = msg.value
		this.filteredProps[this.cursor].Values[msg.scope] = msg.value

		this.valueInput.Reset()

		cmds := tea.Batch(
			Cmd_DisplayMessage("Changed value of "+msg.name+" to "+msg.value),
			Cmd_AddLog(msg.command),
		)

		return this, cmds

	case Msg_DisplayMessage:
		this.message = msg.message
		return this, nil

	case Msg_GetGitConfigProps:
		this.props = msg.data
		this.filteredProps = msg.data
		return this, Cmd_AddLog(msg.command)

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, SearchKeymap.FilterOnlyWithValue):
			this.onlyWithValue = !this.onlyWithValue
			this.filteredProps = git.FilterGitConfigProps(this.props, this.searchInput.Value(), this.onlyWithValue)
			return this, nil

		case key.Matches(msg, SearchKeymap.Quit):
			return this, Cmd_Quit()

		case key.Matches(msg, SearchKeymap.ChangeMode):
			return this, Cmd_SwitchSubmodel(APP_MODEL_LOGS)

		case key.Matches(msg, SearchKeymap.Down):
			if this.cursor < len(this.filteredProps)-1 {
				this.cursor++
			}
			if this.cursor >= this.listStart+10 {
				this.listStart = this.cursor - 10
			}

		case key.Matches(msg, SearchKeymap.Up):
			if this.cursor > 0 {
				this.cursor--
			}
			if this.cursor < this.listStart {
				this.listStart = this.cursor
			}

		case key.Matches(msg, SearchKeymap.PageDown):
			this.cursor += 10
			this.listStart += 10

			if this.cursor > len(this.props) {
				this.cursor = len(this.props) - 9
			}
			if this.listStart > len(this.props) {
				this.listStart = len(this.props) - 9
			}

		case key.Matches(msg, SearchKeymap.PageUp):
			this.cursor -= 10
			this.listStart -= 10

			if this.cursor < 0 {
				this.cursor = 0
			}
			if this.listStart < 0 {
				this.listStart = 0
			}

		case key.Matches(msg, SearchKeymap.NextScope):
			switch this.scope {
			case GIT_SCOPE_LOCAL:
				this.scope = GIT_SCOPE_GLOBAL
			case GIT_SCOPE_GLOBAL:
				this.scope = GIT_SCOPE_SYSTEM
			case GIT_SCOPE_SYSTEM:
				this.scope = GIT_SCOPE_LOCAL
			default:
				this.scope = GIT_SCOPE_LOCAL
			}

		case key.Matches(msg, SearchKeymap.PreviousScope):
			switch this.scope {
			case GIT_SCOPE_LOCAL:
				this.scope = GIT_SCOPE_SYSTEM
			case GIT_SCOPE_GLOBAL:
				this.scope = GIT_SCOPE_LOCAL
			case GIT_SCOPE_SYSTEM:
				this.scope = GIT_SCOPE_GLOBAL
			default:
				this.scope = GIT_SCOPE_LOCAL
			}

		case key.Matches(msg, SearchKeymap.Cancel):
			this.isEditing = false
			return this, nil

		case key.Matches(msg, SearchKeymap.Help):
			this.help.ShowAll = !this.help.ShowAll
			return this, nil

		case key.Matches(msg, SearchKeymap.Confirm):
			if this.isEditing {
				this.isEditing = false
				this.searchInput.Focus()
				this.valueInput.Blur()

				name := this.filteredProps[this.cursor].GetName()
				value := this.valueInput.Value()

				return this, Cmd_GitConfigSet(name, value, this.scope)
			} else {
				this.isEditing = true
				this.searchInput.Blur()
				this.valueInput.Focus()

				oldValue, ok := this.filteredProps[this.cursor].Values[getScopeFromGitScope(this.scope)]
				if ok {
					this.valueInput.SetValue(oldValue)
				} else {
					this.valueInput.SetValue("")
				}

				return this, nil
			}

		}
	}

	if this.isEditing {
		model, cmd := this.valueInput.Update(msg)
		this.valueInput = model
		return this, cmd

	} else {
		prevValue := this.searchInput.Value()
		model, cmd := this.searchInput.Update(msg)
		this.searchInput = model

		if prevValue != this.searchInput.Value() {
			this.filteredProps = git.FilterGitConfigProps(this.props, this.searchInput.Value(), this.onlyWithValue)
			return this, cmd
		}

		return this, cmd

	}
}

func (this MainModel) View() string {
	output := ""

	if this.isEditing {
		output += this.valueInput.View() + "\n"
	} else {
		output += this.searchInput.View() + "\n"
	}

	output += renderGap(80) + "\n"

	for index, prop := range this.filteredProps {
		if index < this.listStart || index > this.listStart+10 {
			continue
		}

		output += renderProp(prop.GetName(), getValueFromScope(prop, this.scope), getColorFromScope(this.scope), index == this.cursor)
	}

	scopeStyle := lipgloss.NewStyle().Foreground(getColorFromScope(this.scope))
	scopeInfo := "--- " + scopeStyle.Render("Scope: "+string(this.scope)) + " ---"
	output += scopeInfo + renderGap(80-len(this.scope)-15) + "\n"
	output += "Last message: " + this.message + "\n"
	output += CenterStyle.Render(this.help.View(SearchKeymap))

	return output
}

func (this MainModel) Help() string {
	return this.help.View(SearchKeymap)
}

// Helpers

func CreateNewMainModel() MainModel {
	searchInput := textinput.New()
	searchInput.Width = 80
	searchInput.Prompt = "Search: "

	valueInput := textinput.New()
	valueInput.Width = 80
	valueInput.Prompt = "Value: "

	help := help.New()
	help.Width = 80

	searchInput.Focus()

	return MainModel{
		searchInput: searchInput,
		valueInput:  valueInput,
		help:        help,
		scope:       GIT_SCOPE_LOCAL,
	}
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

// Commands

type Msg_GetGitConfigProps struct {
	data    []git.GitConfigProp
	command string
}

func Cmd_GetGitConfigEntries() tea.Cmd {
	return func() tea.Msg {
		return Msg_GetGitConfigProps{
			data:    git.GetConfigProps(),
			command: "git config list --show-scope",
		}
	}
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
	command string
	result  bool
	message string
	name    string
	value   string
	scope   string
}

func Cmd_GitConfigSet(name string, value string, gitScope GitScope) tea.Cmd {
	return func() tea.Msg {
		scope := getScopeFromGitScope(gitScope)
		command := fmt.Sprintf("git config set --%s %s %s", scope, name, value)
		content, err := utils.ExecuteCommand("git", "config", "set", "--"+scope, name, value)

		if err != nil {
			return Msg_GitConfigSetResult{
				message: err.Error(),
				result:  false,
				name:    name,
				value:   value,
				scope:   scope,
				command: command,
			}
		}

		return Msg_GitConfigSetResult{
			result:  true,
			message: content,
			name:    name,
			value:   value,
			scope:   scope,
			command: command,
		}
	}
}
