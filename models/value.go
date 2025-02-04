package models

import (
	"github.com/Gabryjiel/git_config_manager/git"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ValueModel struct {
	selectedProp *git.GitConfigProp
	props        []git.GitConfigProp
	input        textinput.Model
	help         help.Model
	keymap       ValueKeymap
	listModel    ListModel[git.GitConfigProp]
}

func CreateNewValueModel() ValueModel {
	listModel := CreateNewListModel(renderInactiveRow)
	helpModel := help.New()

	return ValueModel{keymap: keymap, listModel: listModel, help: helpModel}
}

func (this *ValueModel) Init() tea.Cmd {
	return nil
}

func (this *ValueModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, this.keymap.Cancel):
			this.selectedProp = nil
			return this, nil
		}
	default:
		model, cmd := this.input.Update(msg)
		this.input = model
		return this, cmd
	}

	return this, nil
}

func (this *ValueModel) View() string {
	output := ""

	output += renderEasyHeader(0) + "\n"
	output += this.input.View() + "\n"
	output += renderGap(80) + "\n"

	output += this.listModel.View()

	output += renderGap(80) + "\n"

	output += CenterStyle.Render(this.help.View(this.keymap))

	return output
}

type ValueKeymap struct {
	Cancel key.Binding
}

var keymap = ValueKeymap{
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("ESC", "Cancel"),
	),
}

func (this ValueKeymap) ShortHelp() []key.Binding {
	return []key.Binding{this.Cancel}
}

func (this ValueKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{this.Cancel}}
}
