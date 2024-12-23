package models

import (
	"slices"

	"github.com/Gabryjiel/git_config_manager/git"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewSearchModel() SearchModel {
	systemOptions := git.GetGitConfigByLevel("system")
	globalOptions := git.GetGitConfigByLevel("global")
	localOptions := git.GetGitConfigByLevel("local")

	options := git.CreateConfigMap(systemOptions, globalOptions, localOptions)

	slices.SortFunc(options, func(a, b git.GitConfigProp) int {
		aName := a.GetName()
		bName := b.GetName()

		if aName > bName {
			return 1
		} else if aName < bName {
			return -1
		} else {
			return 0
		}
	})

	return SearchModel{
		cursor:          0,
		allOptions:      options,
		filteredOptions: options,
	}
}

type SearchModel struct {
	tea.Model
	input           TextField
	allOptions      []git.GitConfigProp
	filteredOptions []git.GitConfigProp
	cursor          MenuCursor
}

func (this SearchModel) Init() tea.Cmd {
	return nil
}

func (this SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return this, ExitProgram()
		case tea.KeyCtrlW:
			this.input.removeLastWord()
			this.filteredOptions = git.FilterGitConfigProps(this.allOptions, this.input.Value)
		case tea.KeyCtrlN:
			fallthrough
		case tea.KeyDown:
			this.cursor.moveDown(len(this.filteredOptions))
		case tea.KeyUp:
			fallthrough
		case tea.KeyCtrlP:
			this.cursor.moveUp()
		case tea.KeyBackspace:
			this.input.removeLastCharacter()
			this.filteredOptions = git.FilterGitConfigProps(this.allOptions, this.input.Value)
		case tea.KeyEnter:
			if len(this.filteredOptions) > 0 {
				return this, tea.Sequence(SwitchModel(MODEL_SCOPE), ChooseScopeModelProp(this.filteredOptions[this.cursor]))
			}
		default:
			key := msg.String()

			if len(key) == 1 {
				this.input.Value += key
				this.filteredOptions = git.FilterGitConfigProps(this.allOptions, this.input.Value)

				if int(this.cursor) > len(this.filteredOptions)-1 {
					this.cursor = MenuCursor(len(this.filteredOptions) - 1)
				}
			}
		}
	}

	return this, nil
}

func (this SearchModel) View() string {
	output := renderHeader()
	output += arrowCinStyle.Render(" Search: ")
	output += this.input.Value
	output += lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15)).Render("█")
	output += "\n\n"

	for index, option := range this.filteredOptions {
		output += RenderSearchItem(option, index, int(this.cursor))
	}

	return output
}

func RenderSearchItem(item git.GitConfigProp, itemIndex, cursorIndex int) string {
	defaultStyle := lipgloss.NewStyle()

	if itemIndex == cursorIndex {
		defaultStyle = defaultStyle.Background(lipgloss.ANSIColor(8))
	}

	name := defaultStyle.Width(40).Align(lipgloss.Left).PaddingLeft(1).Render(item.Section + "." + item.Key)

	value := ""
	if len(item.Values.Local) > 0 {
		value = localValueStyle.Background(defaultStyle.GetBackground()).PaddingLeft(1).Render(item.Values.Local)
	} else if len(item.Values.Global) > 0 {
		value = globalValueStyle.Background(defaultStyle.GetBackground()).PaddingLeft(1).Render(item.Values.Global)
	} else if len(item.Values.System) > 0 {
		value = systemValueStyle.Background(defaultStyle.GetBackground()).PaddingLeft(1).Render(item.Values.System)
	}

	value = defaultStyle.Align(lipgloss.Right).Width(40).PaddingRight(1).Render(value)

	return name + value + "\n"
}
