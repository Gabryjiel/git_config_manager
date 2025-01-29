package git

import (
	"log"
	"maps"
	"slices"
	"strings"

	"github.com/Gabryjiel/git_config_manager/utils"
)

const (
	TYPE_STRING    = "string"
	SOURCE__LOCAL  = "local"
	SOURCE__GLOBAL = "global"
	SOURCE__SYSTEM = "system"
)

type ValueScope int

const (
	SCOPE_LOCAL ValueScope = iota
	SCOPE_GLOBAL
	SCOPE_SYSTEM
)

type GitConfigEntry struct {
	Name  string
	Value string
}

func GetConfigProps() []GitConfigProp {
	contents, err := utils.ExecuteCommand("git", "config", "list", "--show-scope")

	if err != nil {
		log.Fatalln("Could not find config list")
	}

	configMap := make(map[string]GitConfigProp)

	lines := strings.Split(contents, "\n")
	for _, line := range lines {
		split := strings.Split(line, "\t")
		scope := split[0]
		entry := split[1]

		splitEntry := strings.SplitN(entry, "=", 2)
		sectionKey := splitEntry[0]
		value := splitEntry[1]

		splitSectionKey := strings.SplitN(sectionKey, ".", 2)
		section := splitSectionKey[0]
		key := splitSectionKey[1]

		_, alreadyExists := configMap[sectionKey]
		if alreadyExists {
			configMap[sectionKey].Values[scope] = value
		} else {
			valuesMap := make(map[string]string)
			valuesMap[scope] = value

			configMap[sectionKey] = GitConfigProp{
				Section: section,
				Key:     key,
				Type:    GIT_CONFIG_PROP__TYPE_STRING,
				Values:  valuesMap,
			}
		}
	}

	configSlice := slices.Collect(maps.Values(configMap))
	slices.SortFunc(configSlice, func(a, b GitConfigProp) int {
		aName := a.GetName()
		bName := b.GetName()

		if aName > bName {
			return 1
		} else if aName < bName {
			return -1
		} else {
			return 0
		}
	})
	return configSlice
}

func GetGitVersion() string {
	result, err := utils.ExecuteCommand("git", "version")
	if err != nil {
		return "git version ?.??.?"
	}

	return result
}

func SetConfigProp(scope ValueScope, key, value string) bool {
	cmdScope := "local"

	switch scope {
	case SCOPE_LOCAL:
		cmdScope = SOURCE__LOCAL
	case SCOPE_GLOBAL:
		cmdScope = SOURCE__GLOBAL
	case SCOPE_SYSTEM:
		cmdScope = SOURCE__SYSTEM
	}

	_, err := utils.ExecuteCommand("git", "config", "--add", "--"+cmdScope, key, value)
	return err != nil
}

func GetConfigLabels() []string {
	labels, err := utils.ExecuteCommand("git", "help", "-c")

	if err != nil {
		return nil
	}

	return strings.Split(labels, "\n")
}
