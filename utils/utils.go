package utils

import (
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
