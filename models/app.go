package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

type AppModel struct {
	tea.Model
	IsExiting         bool
	models            []tea.Model
	CurrentModelIndex int
}

func NewAppModel() AppModel {
	return AppModel{
		IsExiting: false,
		models: []tea.Model{
			NewSearchModel(),
			NewScopeModel(),
		},
		CurrentModelIndex: 0,
	}
}

func (this AppModel) Init() tea.Cmd {
	return nil
}

func (this AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentModel := this.models[this.CurrentModelIndex]
	_, cmd := currentModel.Update(msg)

	return currentModel, cmd
}

func (this AppModel) View() string {
	if this.IsExiting {
		return ""
	}

	return this.models[this.CurrentModelIndex].View()
}
