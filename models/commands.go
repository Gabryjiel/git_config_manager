package models

import (
	"github.com/Gabryjiel/git_config_manager/git"
	tea "github.com/charmbracelet/bubbletea"
)

func SwitchModel(modelId ModelIndex) tea.Cmd {
	return func() tea.Msg {
		return modelId
	}
}

type ExitCode int

func ExitProgram() tea.Cmd {
	return func() tea.Msg {
		return ExitCode(0)
	}
}

func ChooseScopeModelProp(prop git.GitConfigProp) tea.Cmd {
	return func() tea.Msg {
		return prop
	}
}
