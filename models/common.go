package models

import (
	"github.com/Gabryjiel/git_config_manager/git"
	"github.com/charmbracelet/lipgloss"
)

var (
	HeaderStyle = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1))
	CenterStyle = lipgloss.NewStyle().Width(80).AlignHorizontal(lipgloss.Center)
)

func getValueFromScope(prop git.GitConfigProp, scope GitScope) string {
	switch scope {
	case GIT_SCOPE_LOCAL:
		return prop.Values["local"]
	case GIT_SCOPE_GLOBAL:
		return prop.Values["global"]
	case GIT_SCOPE_SYSTEM:
		return prop.Values["system"]
	}

	return ""
}

func getColorFromScope(scope GitScope) lipgloss.ANSIColor {
	switch scope {
	case GIT_SCOPE_LOCAL:
		return lipgloss.ANSIColor(10)
	case GIT_SCOPE_GLOBAL:
		return lipgloss.ANSIColor(11)
	case GIT_SCOPE_SYSTEM:
		return lipgloss.ANSIColor(12)
	default:
		return lipgloss.ANSIColor(1)
	}
}

func getScopeFromGitScope(gitScope GitScope) string {
	switch gitScope {
	case GIT_SCOPE_LOCAL:
		return "local"
	case GIT_SCOPE_GLOBAL:
		return "global"
	case GIT_SCOPE_SYSTEM:
		return "system"
	default:
		return "unknown"
	}
}

func renderGap(length int) string {
	result := ""

	for i := 0; i < length; i++ {
		result += "-"
	}

	return CenterStyle.Render(result)
}

func renderEasyHeader() string {
	result := ""
	result += lipgloss.NewStyle().
		Width(80).
		AlignHorizontal(lipgloss.Center).
		Render("--- gcm v0.0.1 --- " + git.GetGitVersion() + " ---")

	return result
}

func renderHeaderScope(gitScope GitScope) string {
	return lipgloss.NewStyle().Foreground(getColorFromScope(gitScope)).Render(getScopeFromGitScope(gitScope))
}
