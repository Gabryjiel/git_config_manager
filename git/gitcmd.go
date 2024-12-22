package git

import (
	"log"
	"strings"

	"github.com/Gabryjiel/git_config_manager/utils"
)

const (
	TYPE_STRING    = "string"
	SOURCE__LOCAL  = "local"
	SOURCE__GLOBAL = "global"
	SOURCE__SYSTEM = "system"
)

type GitConfigEntry struct {
	Name  string
	Value string
}

func GetGitConfigByLevel(source string) []GitConfigProp {
	cmdOutputStr, err := utils.ExecuteCommand("git", "config", "--list", "--"+source)

	if err != nil {
		log.Println("Failed", err)
		return nil
	}

	entries := strings.Split(cmdOutputStr, "\n")
	result := make([]GitConfigProp, len(entries))

	for i, entryStr := range entries {
		split := strings.SplitN(entryStr, "=", 2)
		location := strings.SplitN(split[0], ".", 2)

		result[i].Type = GIT_CONFIG_PROP__TYPE_STRING
		result[i].Section = location[0]
		result[i].Key = location[1]
		result[i].Values = GitConfigEntryValues{}

		switch source {
		case SOURCE__GLOBAL:
			result[i].Values.Global = split[1]
		case SOURCE__SYSTEM:
			result[i].Values.System = split[1]
		case SOURCE__LOCAL:
			result[i].Values.Local = split[1]
		}
	}

	return result
}

func GetGitVersion() string {
	result, err := utils.ExecuteCommand("git", "version")
	if err != nil {
		return "git version ?.??.?"
	}

	return result
}
