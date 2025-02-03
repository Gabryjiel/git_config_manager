package models

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func NewLogsModel() logsModel {
	help := help.New()

	return logsModel{help: help}
}

type logsModel struct {
	logs   []string
	cursor MenuCursor
	help   help.Model
}

func (this *logsModel) Init() tea.Cmd {
	return nil
}

func (this *logsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, LogsModelKeyMap.Up):
			this.cursor.moveUp(1)

		case key.Matches(msg, LogsModelKeyMap.Down):
			this.cursor.moveDown(len(this.logs), 1)

		case key.Matches(msg, LogsModelKeyMap.PageUp):
			this.cursor.moveUp(10)

		case key.Matches(msg, LogsModelKeyMap.PageDown):
			this.cursor.moveDown(len(this.logs), 10)

		case key.Matches(msg, LogsModelKeyMap.ChangeMode):
			return this, Cmd_ChangeMode(APP_MODE_SEARCH)

		case key.Matches(msg, LogsModelKeyMap.Help):
			this.help.ShowAll = !this.help.ShowAll

		case key.Matches(msg, LogsModelKeyMap.Quit):
			return this, Cmd_Quit()
		}
	}

	return this, nil
}

func (this *logsModel) View() string {
	output := ""

	output += renderEasyHeader(0) + "\n"
	output += renderGap(80) + "\n"

	for _, log := range this.logs {
		output += log + "\n"
	}

	output += renderGap(80) + "\n"
	output += CenterStyle.Render(this.help.View(LogsModelKeyMap))

	return output
}

func (this *logsModel) PushLog(log string) {
	this.logs = append(this.logs, log)
}

type logsModelKeymap struct {
	Up         key.Binding
	Down       key.Binding
	PageDown   key.Binding
	PageUp     key.Binding
	Help       key.Binding
	ChangeMode key.Binding
	Quit       key.Binding
}

var LogsModelKeyMap = logsModelKeymap{
	Up: key.NewBinding(
		key.WithKeys("up", "ctrl+p"),
		key.WithHelp("↑/<C-p>", "Move cursor up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "ctrl+n"),
		key.WithHelp("↓/<C-n>", "Move cursor down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("PgUp", "Go up 10"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("PgDn", "Go down 10"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	ChangeMode: key.NewBinding(
		key.WithKeys("ctrl+up"),
		key.WithHelp("<C-↑>", "Switch view"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("<C-c>", "quit"),
	),
}

func (k logsModelKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k logsModelKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PageUp, k.PageDown, k.Help, k.ChangeMode},
	}
}
