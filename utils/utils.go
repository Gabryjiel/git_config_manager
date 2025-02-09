package utils

import (
	"errors"
	"os/exec"
	"slices"
	"strings"
)

func ExecuteCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name)
	cmd.Args = slices.Concat(cmd.Args, args)
	cmdOutput, err := cmd.Output()

	if err != nil {
		return "", err
	}

	cmdOutputStr := strings.TrimSpace(string(cmdOutput))

	return cmdOutputStr, nil
}

func ExecuteSimpleCommand(command string) (string, error) {
	splitted := strings.Split(command, " ")
	if len(splitted) < 1 {
		return "", errors.New("No command found")
	}

	return ExecuteCommand(splitted[0], splitted[1:]...)
}
