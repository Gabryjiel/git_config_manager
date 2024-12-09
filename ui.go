package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewModel struct {
	cursor          int
	options         []GitConfigEntry
	filteredOptions []GitConfigEntry
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

var headerStyle lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("01"))
var selectedStyle lipgloss.Style = lipgloss.NewStyle().Background(lipgloss.Color("00"))

func (model ViewModel) View() string {
	if model.isExiting {
		return ""
	}

	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))

	output := style.Render(" >> ")
	output += model.searchPhrase
	output += lipgloss.NewStyle().Background(lipgloss.Color("11")).Render(" ")
	output += "\n"
	output += headerStyle.Render("--- Git Config Manager v0.0.1 --- " + GetGitVersion() + " ---")
	output += "\n"

	for index, option := range model.filteredOptions {
		if index == model.cursor {
			output += " "
			output += selectedStyle.Render(option.Section + "." + option.Key)
		} else {
			output += " "
			output += option.Section + "." + option.Key
		}

		output += "\t" + option.Value["local"] + " / " + option.Value["global"] + " / " + option.Value["system"]
		output += "\n"
	}

	return output
}

func (model *ViewModel) MoveCursorDown() {
	if model.cursor < len(model.options)-1 {
		model.cursor++
	}

}

func (model *ViewModel) MoveCursorUp() {
	if model.cursor > 0 {
		model.cursor--
	}
}

func (model *ViewModel) FilterOptions() {
	model.filteredOptions = filterOptions(model.options, model.searchPhrase)
}

func createInitialModel() ViewModel {
	options := GetGitConfigByLevel("local")

	return ViewModel{
		cursor:          0,
		options:         options,
		filteredOptions: options,
		isExiting:       false,
	}
}
