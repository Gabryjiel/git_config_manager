package models

import tea "github.com/charmbracelet/bubbletea"

type ListModel[TElement any] struct {
	cursor     MenuCursor
	upperBound uint
	items      []TElement
	renderRow  func(TElement) string
}

func CreateNewListModel[TElement any](renderRow func(TElement) string) ListModel[TElement] {
	return ListModel[TElement]{
		renderRow: renderRow,
	}
}

func (this *ListModel[T]) Init() tea.Cmd {
	return nil
}

func (this *ListModel[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return this, nil
}

func (this *ListModel[T]) View() string {
	output := ""

	for _, item := range this.items {
		output += this.renderRow(item) + "\n"
	}

	return output
}

//
