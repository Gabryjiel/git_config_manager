package main

import (
	"fmt"
	"os"

	"github.com/Gabryjiel/git_config_manager/models"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(models.CreateNewMainModel())
	if _, err := program.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
