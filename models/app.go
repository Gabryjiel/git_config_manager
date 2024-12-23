package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ModelIndex int

const (
	MODEL_SEARCH ModelIndex = iota
	MODEL_SCOPE
)

type AppModel struct {
	tea.Model
	IsExiting         bool
	models            []tea.Model
	CurrentModelIndex ModelIndex
}

func NewAppModel() AppModel {
	return AppModel{
		IsExiting: false,
		models: []tea.Model{
			NewSearchModel(),
			NewScopeModel(),
		},
		CurrentModelIndex: MODEL_SEARCH,
	}
}

func (this AppModel) Init() tea.Cmd {
	return nil
}

func (this AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ModelIndex:
		this.CurrentModelIndex = msg
	case ExitCode:
		this.IsExiting = true
		return this, tea.Quit
	}

	currentModel := this.models[this.CurrentModelIndex]
	newModel, cmd := currentModel.Update(msg)
	this.models[this.CurrentModelIndex] = newModel

	return this, cmd
}

func (this AppModel) View() string {
	if this.IsExiting {
		return ""
	}

	return this.models[this.CurrentModelIndex].View()
}
