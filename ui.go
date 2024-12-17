package main

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	headerStyle      = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1))
	arrowCinStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(14))
	localValueStyle  = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(10))
	globalValueStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(11))
	systemValueStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(12))
)

type ViewModel struct {
	cursor          int
	allOptions      []GitConfigProp
	filteredOptions []GitConfigProp
	searchPhrase    string
	isExiting       bool
}

func (model ViewModel) Init() tea.Cmd {
	return nil
}

func (model ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			model.isExiting = true
			return model, tea.Quit
		case tea.KeyCtrlW:
			lastIndex := strings.LastIndex(model.searchPhrase, " ")
			if lastIndex == -1 {
				model.searchPhrase = ""
			} else {
				model.searchPhrase = model.searchPhrase[0:lastIndex]
			}
			model.FilterOptions()
		case tea.KeyCtrlN:
			model.MoveCursorDown()
		case tea.KeyCtrlP:
			model.MoveCursorUp()
		case tea.KeyBackspace:
			if len(model.searchPhrase) > 0 {
				model.searchPhrase = model.searchPhrase[:len(model.searchPhrase)-1]
				model.FilterOptions()
			}
		default:
			key := msg.String()

			if len(key) == 1 {
				model.searchPhrase += key
				model.FilterOptions()

				if model.cursor > len(model.filteredOptions)-1 {
					model.cursor = len(model.filteredOptions) - 1
				}
			}
		}
	}

	return model, nil
}

func (model ViewModel) View() string {
	if model.isExiting {
		return ""
	}

	output := arrowCinStyle.Render(" >> ")
	output += model.searchPhrase
	output += lipgloss.NewStyle().Background(lipgloss.ANSIColor(11)).Render(" ")
	output += "\n"
	output += headerStyle.Render("--- GCM v0.0.1 --- " + GetGitVersion() + " --- ")
	output += localValueStyle.Render(" L ")
	output += globalValueStyle.Render(" G ")
	output += systemValueStyle.Render(" S ")
	output += headerStyle.Render("---")
	output += "\n"

	for index, option := range model.filteredOptions {
		output += renderItem(option, index, model.cursor)
		output += "\n"
	}

	return output
}

func (model *ViewModel) MoveCursorDown() {
	if model.cursor < len(model.allOptions)-1 {
		model.cursor++
	}

}

func (model *ViewModel) MoveCursorUp() {
	if model.cursor > 0 {
		model.cursor--
	}
}

func (model *ViewModel) FilterOptions() {
	if len(model.searchPhrase) == 0 {
		model.filteredOptions = model.allOptions
	}

	filteredOptions := make([]GitConfigProp, 0)

	for _, option := range model.allOptions {
		if strings.Contains(option.toString(), model.searchPhrase) {
			filteredOptions = append(filteredOptions, option)
		}
	}

	model.filteredOptions = filteredOptions
}

func createInitialModel() ViewModel {
	systemOptions := GetGitConfigByLevel("system")
	globalOptions := GetGitConfigByLevel("global")
	localOptions := GetGitConfigByLevel("local")

	options := CreateConfigMap(systemOptions, globalOptions, localOptions)

	slices.SortFunc(options, func(a, b GitConfigProp) int {
		aName := a.getName()
		bName := b.getName()

		if aName > bName {
			return 1
		} else if aName < bName {
			return -1
		} else {
			return 0
		}
	})

	return ViewModel{
		cursor:          0,
		allOptions:      options,
		filteredOptions: options,
		isExiting:       false,
	}
}

func renderItem(item GitConfigProp, itemIndex, cursorIndex int) string {
	defaultStyle := lipgloss.NewStyle()

	if itemIndex == cursorIndex {
		defaultStyle = defaultStyle.Background(lipgloss.ANSIColor(8)).Bold(true)
	}

	name := defaultStyle.Width(40).Align(lipgloss.Left).PaddingLeft(1).Render(item.Section + "." + item.Key)

	values := make([]string, 0)
	if len(item.Values.Local) > 0 {
		values = append(values, localValueStyle.Background(defaultStyle.GetBackground()).PaddingLeft(1).Render(item.Values.Local))
	}
	if len(item.Values.Global) > 0 {
		values = append(values, globalValueStyle.Background(defaultStyle.GetBackground()).PaddingLeft(1).Render(item.Values.Global))
	}
	if len(item.Values.System) > 0 {
		values = append(values, systemValueStyle.Background(defaultStyle.GetBackground()).PaddingLeft(1).Render(item.Values.System))
	}

	value := strings.Join(values, "")
	value = defaultStyle.Align(lipgloss.Right).Width(40).PaddingRight(1).Render(value)

	return name + value
}
