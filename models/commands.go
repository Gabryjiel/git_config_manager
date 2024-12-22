package models

import tea "github.com/charmbracelet/bubbletea"

func SwitchModel(modelId int) tea.Cmd {
	return func() tea.Msg {
		return modelId
	}
}
