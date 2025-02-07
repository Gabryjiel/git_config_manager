package models

import tea "github.com/charmbracelet/bubbletea"

func NewAppModel() *AppModel {
	model := AppModel{
		isExiting: false,
		cursor:    APP_MODEL_LIST,
		models: [2]tea.Model{
			CreateNewMainModel(),
			NewLogsModel(),
		},
	}

	return &model
}

type Submodel interface {
	tea.Model
	Help() string
}

type AppSubmodelId int

const (
	APP_MODEL_LIST AppSubmodelId = iota
	APP_MODEL_LOGS
)

type AppModel struct {
	isExiting bool
	cursor    AppSubmodelId
	models    [2]tea.Model
}

func (this AppModel) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, submodel := range this.models {
		cmds = append(cmds, submodel.Init())
	}

	return tea.Batch(cmds...)
}

func (this AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			this.isExiting = true
			return this, tea.Quit
		}

	case Msg_SwitchSubmodel:
		this.cursor = msg.submodelId

		return this, nil

	case Msg_NewLog:
		submodel, cmd := this.models[APP_MODEL_LOGS].Update(msg)
		this.models[APP_MODEL_LOGS] = submodel
		return this, cmd

	case Msg_Quit:
		this.isExiting = true
		return this, tea.Quit
	}

	submodel, cmd := this.models[this.cursor].Update(msg)
	this.models[this.cursor] = submodel

	return this, cmd
}

func (this AppModel) View() string {
	output := ""

	if this.isExiting {
		return output
	}

	output += renderEasyHeader() + "\n"
	output += renderGap(80) + "\n"

	output += this.models[this.cursor].View() + "\n"

	return output
}

// Commands

type Msg_SwitchSubmodel struct {
	submodelId AppSubmodelId
}

func Cmd_SwitchSubmodel(submodelId AppSubmodelId) tea.Cmd {
	return func() tea.Msg {
		return Msg_SwitchSubmodel{submodelId}
	}
}

type Msg_Quit struct{}

func Cmd_Quit() tea.Cmd {
	return func() tea.Msg {
		return Msg_Quit{}
	}
}
