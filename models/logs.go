package models

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewLogsModel() LogsModel {
	help := help.New()

	return LogsModel{help: help}
}

type Log struct {
	message   string
	timestamp time.Time
}

type LogsModel struct {
	logs   []Log
	cursor MenuCursor
	help   help.Model
}

func (this LogsModel) Init() tea.Cmd {
	return nil
}

func (this LogsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Msg_NewLog:
		this.logs = append(this.logs, Log(msg))

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
			return this, Cmd_SwitchSubmodel(APP_MODEL_LIST)

		case key.Matches(msg, LogsModelKeyMap.Help):
			this.help.ShowAll = !this.help.ShowAll

		case key.Matches(msg, LogsModelKeyMap.Quit):
			return this, Cmd_Quit()
		}
	}

	return this, nil
}

func (this LogsModel) View() string {
	output := ""

	output += "Logs (git commands)" + "\n"
	output += renderGap(80) + "\n"

	logStyle := lipgloss.NewStyle().Width(80).PaddingLeft(1).PaddingRight(1)
	timeStyle := lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(15))

	for _, log := range this.logs {
		formatted := fmt.Sprintf("[%s] %s", timeStyle.Render(log.timestamp.Format(time.TimeOnly)), log.message)
		output += logStyle.Render(formatted) + "\n"
	}

	if len(this.logs) == 0 {
		output += " No logs \n"
	}

	output += renderGap(80) + "\n"
	output += CenterStyle.Render(this.help.View(LogsModelKeyMap))

	return output
}

// Helpers

func (this *LogsModel) PushLog(log Log) {
	this.logs = append(this.logs, log)
}

type Msg_NewLog Log

func Cmd_AddLog(log string) tea.Cmd {
	return func() tea.Msg {
		return Msg_NewLog{message: log, timestamp: time.Now()}
	}
}

// Keymap

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
