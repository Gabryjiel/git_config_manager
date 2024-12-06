package main

import (
	"fmt"
	"log"
	"maps"
	"os"
	"os/exec"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	program := tea.NewProgram(createInitialModel())
	if _, err := program.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

type ViewModel struct {
	cursor          int
	options         []ConfigEntry
	filteredOptions []ConfigEntry
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
	output += headerStyle.Render("--- Git Config Manager v0.0.1 --- " + getGitVersion() + " ---")
	output += "\n"

	for index, option := range model.filteredOptions {
		if index == model.cursor {
			output += " "
			output += selectedStyle.Render(option.Name)
		} else {
			output += " "
			output += option.Name
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
	optionsMap := getList()
	options := slices.Collect(maps.Values(optionsMap))

	return ViewModel{
		cursor:          0,
		options:         options,
		filteredOptions: options,
		isExiting:       false,
	}
}

func filterOptions(options []ConfigEntry, searchPhrase string) []ConfigEntry {
	if len(searchPhrase) == 0 {
		return options
	}

	filteredOptions := make([]ConfigEntry, 0)

	for _, option := range options {
		if strings.Contains(option.Name, searchPhrase) {
			filteredOptions = append(filteredOptions, option)
		}
	}

	return filteredOptions
}

type ConfigEntry struct {
	Name    string
	Section string
	Key     string
	Type    string
	Value   map[string]string
}

func getGitVersion() string {
	cmd := exec.Command("git", "version")
	result, err := cmd.Output()

	if err != nil {
		return ""
	}

	resultStr := string(result)

	return strings.TrimSpace(resultStr)
}

type ConfigEntryStr struct {
	Key    string
	Value  string
	Source string
}

func parseGitListOutput(source string) []ConfigEntryStr {
	cmdOutputStr, err := executeCommand("git", "config", "--list", "--"+source)

	if err != nil {
		log.Println("Failed", err)
		return nil
	}

	entries := strings.Split(cmdOutputStr, "\n")

	result := make([]ConfigEntryStr, len(entries)-1)

	for index, entryStr := range entries {
		split := strings.Split(entryStr, "=")

		if len(split) != 2 {
			continue
		}

		result[index].Key = split[0]
		result[index].Value = split[1]
		result[index].Source = source
	}

	return result
}

func executeCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name)
	cmd.Args = slices.Concat(cmd.Args, args)
	cmdOutput, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return string(cmdOutput), nil
}

func getList() map[string]ConfigEntry {
	localEntries := parseGitListOutput("local")
	globalEntries := parseGitListOutput("global")
	systemEntries := parseGitListOutput("system")
	entries := slices.Concat(localEntries, globalEntries, systemEntries)

	entryMap := make(map[string]ConfigEntry)

	for _, entry := range entries {
		location := strings.Split(entry.Key, ".")

		if len(location) != 2 {
			continue
		}

		_, ok := entryMap[entry.Key]

		if ok {
			entryMap[entry.Key].Value[entry.Source] = entry.Value
		} else {
			entryMap[entry.Key] = ConfigEntry{
				Name:    entry.Key,
				Section: location[0],
				Key:     location[1],
				Value:   make(map[string]string),
				Type:    "",
			}

			entryMap[entry.Key].Value[entry.Source] = entry.Value
		}
	}

	return entryMap
}
