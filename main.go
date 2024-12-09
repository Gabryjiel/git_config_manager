package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(createInitialModel())
	if _, err := program.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func filterOptions(options []GitConfigEntry, searchPhrase string) []GitConfigEntry {
	if len(searchPhrase) == 0 {
		return options
	}

	filteredOptions := make([]GitConfigEntry, 0)

	for _, option := range options {
		if strings.Contains(option.EntryStr, searchPhrase) {
			filteredOptions = append(filteredOptions, option)
		}
	}

	return filteredOptions
}

func executeCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name)
	cmd.Args = slices.Concat(cmd.Args, args)
	cmdOutput, err := cmd.Output()

	if err != nil {
		return "", err
	}

	cmdOutputStr := strings.TrimSpace(string(cmdOutput))

	return cmdOutputStr, nil
}
