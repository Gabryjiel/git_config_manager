package models

import tea "github.com/charmbracelet/bubbletea"

type AppModel struct {
	isExiting bool
	models    [2]tea.Model
	cursor    AppSubmodelId
}

type AppSubmodelId int

const (
	MODEL_ID_SEARCH AppSubmodelId = iota
	MODEL_ID_SCOPE
)

func (this *AppModel) Init() tea.Cmd {
	this.models[MODEL_ID_SEARCH] = NewSearchModel()
	this.models[MODEL_ID_SCOPE] = NewScopeModel()

	return nil
}

func (this *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			this.isExiting = true
			return this, tea.Quit
		}
	case AppSubmodelId:
		this.cursor = msg
		return this, this.models[this.cursor].Init()
	}

	_, cmd := this.models[this.cursor].Update(msg)
	return this, cmd
}

func (this *AppModel) View() string {
	if this.isExiting {
		return ""
	}

	return this.models[this.cursor].View()
}

// Commands

func SwitchSubmodel(submodelId AppSubmodelId) tea.Cmd {
	return func() tea.Msg {
		return submodelId
	}
}
