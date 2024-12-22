package models

import (
	"strings"

	"github.com/Gabryjiel/git_config_manager/git"
)

type TextField struct {
	Value      string
	IsDisabled bool
}

func (this *TextField) removeLastCharacter() {
	if len(this.Value) > 0 {
		this.Value = this.Value[:len(this.Value)-1]
	}
}

func (this *TextField) removeLastWord() {
	lastIndex := strings.LastIndex(this.Value, " ")
	if lastIndex == -1 {
		this.Value = ""
	} else {
		this.Value = this.Value[0:lastIndex]
	}
}

type MenuCursor int

func (this *MenuCursor) moveDown(max int) {
	if int(*this) < max {
		*this++
	}
}

func (this *MenuCursor) moveUp() {
	if *this > 0 {
		*this--
	}
}

func renderHeader() string {
	output := ""

	output += headerStyle.Render("--- GCM v0.0.1 --- " + git.GetGitVersion() + " --- ")
	output += localValueStyle.Render(" L ")
	output += globalValueStyle.Render(" G ")
	output += systemValueStyle.Render(" S ")
	output += headerStyle.Render("---")
	output += "\n"

	return output
}
