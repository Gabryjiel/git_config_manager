package models

import (
	"github.com/Gabryjiel/git_config_manager/git"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewSearchModel() *SearchModel {
	configProps := git.GetConfigProps()

	return &SearchModel{
		cursor:          0,
		allOptions:      configProps,
		filteredOptions: configProps,
		input:           NewTextField(),
	}
}

type SearchModel struct {
	input           TextInputModel
	allOptions      []git.GitConfigProp
	filteredOptions []git.GitConfigProp
	cursor          MenuCursor
}

func (this *SearchModel) Init() tea.Cmd {
	return nil
}

func (this *SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyDown:
			this.cursor.moveDown(len(this.filteredOptions))
			return this, nil

		case tea.KeyUp:
			this.cursor.moveUp()
			return this, nil

		case tea.KeyEnter:
			if len(this.filteredOptions) > 0 {
				return this, SwitchSubmodel(MODEL_ID_SCOPE)
			}
			return this, nil
		}
	case MsgInputChanged:
		this.filteredOptions = git.FilterGitConfigProps(this.allOptions, this.input.GetValue())
		return this, nil
	}

	_, cmd := this.input.Update(msg)
	return this, cmd
}

func (this *SearchModel) View() string {
	output := renderHeader()
	output += arrowCinStyle.Render(" Search: ")
	output += this.input.View()
	output += "\n\n"

	for index, option := range this.filteredOptions {
		output += RenderSearchItem(option, index, int(this.cursor))
	}

	return output
}

// Helpers

func RenderSearchItem(item git.GitConfigProp, itemIndex, cursorIndex int) string {
	defaultStyle := lipgloss.NewStyle()

	if itemIndex == cursorIndex {
		defaultStyle = defaultStyle.Background(lipgloss.ANSIColor(8))
	}

	name := defaultStyle.Width(40).Align(lipgloss.Left).PaddingLeft(1).Render(item.Section + "." + item.Key)

	value := getItemValue(item, defaultStyle)
	value = defaultStyle.Align(lipgloss.Right).Width(40).PaddingRight(1).Render(value)

	return name + value + "\n"
}

func getItemValue(item git.GitConfigProp, defaultStyle lipgloss.Style) string {
	bgColor := defaultStyle.GetBackground()
	value, ok := item.Values["local"]
	if ok {
		return localValueStyle.Background(bgColor).PaddingLeft(1).Render(value)
	}

	value, ok = item.Values["global"]
	if ok {
		return globalValueStyle.Background(bgColor).PaddingLeft(1).Render(value)
	}

	value, ok = item.Values["system"]
	if ok {
		return systemValueStyle.Background(bgColor).PaddingLeft(1).Render(value)
	}

	return defaultStyle.PaddingLeft(1).Render("")
}

// Commands
