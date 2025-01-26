package models

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewTextField() TextInputModel {
	model := TextInputModel{cursorPosition: 0, value: ""}
	return model
}

type TextInputModel struct {
	cursorPosition int
	value          string
}

func (this *TextInputModel) Init() tea.Cmd {
	return nil
}

func (this *TextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlW:
			this.removeLastWord()
			return this, CmdInputChanged
		case tea.KeyBackspace:
			this.removePreviousCharacter()
			return this, CmdInputChanged
		case tea.KeyLeft:
			this.moveCursor(-1)
		case tea.KeyRight:
			this.moveCursor(1)
		default:
			if len(msg.String()) != 1 {
				break
			}

			this.insertOnCursor(msg.String())
			return this, CmdInputChanged
		}
	}

	return this, nil
}

func (this *TextInputModel) View() string {
	cursorStyle := lipgloss.NewStyle().Background(lipgloss.ANSIColor(1)).Foreground(lipgloss.ANSIColor(0))
	value := this.value + " "

	before := value[:this.cursorPosition]
	after := value[this.cursorPosition+1:]
	cursored := cursorStyle.Render(string(value[this.cursorPosition]))

	return before + cursored + after
}

// Methods

func (this *TextInputModel) removePreviousCharacter() {
	if len(this.value) < 1 || this.cursorPosition == 0 {
		return
	}

	before := this.value[:this.cursorPosition-1]
	after := this.value[this.cursorPosition:]
	this.value = before + after
	this.moveCursor(-1)

}

func (this *TextInputModel) removeLastWord() {
	lastIndex := strings.LastIndex(this.value, " ")

	if lastIndex == -1 {
		this.value = ""
		this.cursorPosition = 0
	} else {
		this.value = this.value[0:lastIndex]
		this.cursorPosition = lastIndex
	}
}

func (this *TextInputModel) moveCursor(length int) {
	result := this.cursorPosition + length

	if result < 0 {
		this.cursorPosition = 0
	} else if result > len(this.value) {
		this.cursorPosition = len(this.value)
	} else {
		this.cursorPosition = result
	}
}

func (this *TextInputModel) insertOnCursor(chars string) {
	before := this.value[:this.cursorPosition]
	after := this.value[this.cursorPosition:]

	this.value = before + chars + after
	this.cursorPosition += len(chars)
}

func (this *TextInputModel) GetValue() string {
	return this.value
}

// Commands

type MsgInputChanged struct{}

func CmdInputChanged() tea.Msg {
	return MsgInputChanged{}
}
