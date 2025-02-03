package models

import "github.com/charmbracelet/bubbles/key"

type searchKeymap struct {
	Up                  key.Binding
	Down                key.Binding
	PageDown            key.Binding
	PageUp              key.Binding
	PreviousScope       key.Binding
	NextScope           key.Binding
	FilterOnlyWithValue key.Binding
	Cancel              key.Binding
	Confirm             key.Binding
	ChangeMode          key.Binding
	Help                key.Binding
	Quit                key.Binding
}

var SearchKeymap = searchKeymap{
	Up: key.NewBinding(
		key.WithKeys("up", "ctrl+p"),
		key.WithHelp("↑/<C-p>", "Move cursor up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "ctrl+n"),
		key.WithHelp("↓/<C-n>", "Move cursor down"),
	),
	PreviousScope: key.NewBinding(
		key.WithKeys("ctrl+left"),
		key.WithHelp("<C-←>", "Previous scope"),
	),
	NextScope: key.NewBinding(
		key.WithKeys("ctrl+right"),
		key.WithHelp("<C-→>", "Next scope"),
	),
	FilterOnlyWithValue: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("Tab", "Only with value"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("ESC", "Cancel"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("<CR>", "Confirm"),
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

func (k searchKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k searchKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.PageUp, k.PageDown, k.Help, k.ChangeMode},
		{k.PreviousScope, k.NextScope, k.FilterOnlyWithValue, k.Cancel, k.Quit},
	}
}
