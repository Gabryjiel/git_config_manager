package models

import (
	"slices"

	"github.com/Gabryjiel/git_config_manager/git"
	"github.com/Gabryjiel/git_config_manager/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GitScope int
type AppMode uint8

const (
	GIT_SCOPE_LOCAL GitScope = iota
	GIT_SCOPE_GLOBAL
	GIT_SCOPE_SYSTEM
)

const (
	APP_MODE_SEARCH AppMode = iota
	APP_MODE_VALUE
	APP_MODE_LOGS
)

type MainModel struct {
	isExiting     bool
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
	mode          AppMode
	logsModel     logsModel
}

func (this *MainModel) Init() tea.Cmd {
	return Cmd_GetGitConfigEntries()
}

func (this *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case Msg_GitConfigSetResult:
		if !msg.result {
			return this, Cmd_DisplayMessage(msg.message)
		}

		this.mode = APP_MODE_SEARCH
		index := slices.IndexFunc(this.props, func(prop git.GitConfigProp) bool {
			return prop.GetName() == msg.name
		})

		this.props[index].Values[msg.scope] = msg.value
		this.filteredProps[this.cursor].Values[msg.scope] = msg.value

		this.valueInput.Reset()

		return this, Cmd_DisplayMessage("Changed value of " + msg.name + " to " + msg.value)

	case Msg_DisplayMessage:
		this.message = msg.message
		return this, nil

	case Msg_GetGitConfigProps:
		this.props = msg.data
		this.filteredProps = msg.data
		return this, nil

	case Msg_Quit:
		this.isExiting = true
		return this, tea.Quit

	case Msg_ChangeMode:
		this.mode = msg.mode

	case Msg_InputChanged:
		this.filteredProps = git.FilterGitConfigProps(this.props, this.searchInput.Value(), this.onlyWithValue)

		if len(this.filteredProps) == 0 {
			this.cursor = 0
		} else if this.cursor >= len(this.filteredProps) {
			this.cursor = len(this.filteredProps) - 1
		}

		return this, nil

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, SearchKeymap.FilterOnlyWithValue):
			this.onlyWithValue = !this.onlyWithValue
			return this, Cmd_InputChanged

		case key.Matches(msg, SearchKeymap.Quit):
			this.isExiting = true
			return this, tea.Quit

		case key.Matches(msg, SearchKeymap.ChangeMode):
			switch {
			case this.mode == APP_MODE_LOGS:
				this.mode = APP_MODE_SEARCH
			case this.mode == APP_MODE_SEARCH:
				this.mode = APP_MODE_LOGS
			}

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
			if this.scope == 2 {
				this.scope = 0
			} else {
				this.scope++
			}

		case key.Matches(msg, SearchKeymap.PreviousScope):
			if this.scope == 0 {
				this.scope = 2
			} else {
				this.scope--
			}

		case key.Matches(msg, SearchKeymap.Cancel):
			this.mode = APP_MODE_SEARCH
			return this, nil

		case key.Matches(msg, SearchKeymap.Help):
			this.help.ShowAll = !this.help.ShowAll
			return this, nil

		case key.Matches(msg, SearchKeymap.Confirm):
			switch this.mode {
			case APP_MODE_SEARCH:
				this.mode = APP_MODE_VALUE
				this.searchInput.Blur()
				this.valueInput.Focus()

				oldValue, ok := this.filteredProps[this.cursor].Values[getScopeFromGitScope(this.scope)]
				if ok {
					this.valueInput.SetValue(oldValue)
				} else {
					this.valueInput.SetValue("")
				}

				return this, nil

			case APP_MODE_VALUE:
				name := this.filteredProps[this.cursor].GetName()
				value := this.valueInput.Value()

				this.searchInput.Focus()
				this.valueInput.Blur()

				return this, Cmd_GitConfigSet(name, value, this.scope)
			}
		}
	}

	switch this.mode {
	case APP_MODE_SEARCH:
		prevValue := this.searchInput.Value()
		model, cmd := this.searchInput.Update(msg)
		this.searchInput = model

		if prevValue != this.searchInput.Value() {
			return this, tea.Batch(Cmd_InputChanged, cmd)
		}

		return this, cmd

	case APP_MODE_VALUE:
		model, cmd := this.valueInput.Update(msg)
		this.valueInput = model
		return this, cmd
	}

	return this, nil
}

func (this *MainModel) View() string {
	output := ""

	if this.isExiting {
		return ""
	} else if this.mode == APP_MODE_LOGS {
		output += this.logsModel.View()
	} else {
		output += renderEasyHeader(this.scope) + "\n"

		if this.mode == APP_MODE_VALUE {
			output += this.valueInput.View() + "\n"
		} else if this.mode == APP_MODE_SEARCH {
			output += this.searchInput.View() + "\n"
		}

		output += renderGap(80) + "\n"

		for index, prop := range this.filteredProps {
			if index < this.listStart || index > this.listStart+10 {
				continue
			}

			output += renderProp(prop.GetName(), getValueFromScope(prop, this.scope), getColorFromScope(this.scope), index == this.cursor)
		}

		output += renderGap(80) + "\n"
		output += "Last message: " + this.message + "\n"
		output += CenterStyle.Render(this.help.View(SearchKeymap))
	}

	return output
}

// Helpers

func CreateNewMainModel() *MainModel {
	searchInput := textinput.New()
	searchInput.Width = 80
	searchInput.Prompt = "Search: "

	valueInput := textinput.New()
	valueInput.Width = 80
	valueInput.Prompt = "Value: "

	help := help.New()
	help.Width = 80

	searchInput.Focus()

	return &MainModel{
		searchInput: searchInput,
		valueInput:  valueInput,
		help:        help,
		logsModel:   NewLogsModel(),
	}
}

func renderGap(length int) string {
	result := ""

	for i := 0; i < length; i++ {
		result += "-"
	}

	return CenterStyle.Render(result)
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

type Msg_GetGitConfigProps struct {
	data []git.GitConfigProp
}

func Cmd_GetGitConfigEntries() tea.Cmd {
	return func() tea.Msg {
		return Msg_GetGitConfigProps{data: git.GetConfigProps()}
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

type Msg_ChangeMode struct {
	mode AppMode
}

func Cmd_ChangeMode(mode AppMode) tea.Cmd {
	return func() tea.Msg {
		return Msg_ChangeMode{mode}
	}
}

type Msg_Quit struct{}

func Cmd_Quit() tea.Cmd {
	return func() tea.Msg {
		return Msg_Quit{}
	}
}
