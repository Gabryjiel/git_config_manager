package models

import tea "github.com/charmbracelet/bubbletea"

type ListModel struct{}

func (this *ListModel) Init() tea.Cmd {
	return nil
}

func (this *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return this, nil
}

func (this *ListModel) View() string {
	return ""
}
