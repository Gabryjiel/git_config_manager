package git

import (
	"errors"
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
	Scope string
}

type GitConfigMap map[string]GitConfigProp

func ParseScopedGitConfigList(contents string) []GitConfigEntry {
	lines := strings.Split(contents, "\n")
	entries := make([]GitConfigEntry, len(lines))
	for index, line := range lines {
		entry, err := ParseScopedGitEntry(line)
		if err != nil {
			continue
		}

		entries[index] = entry
	}

	return entries
}

func ParseScopedGitEntry(toParse string) (GitConfigEntry, error) {
	parsedEntry := GitConfigEntry{}

	splitted := strings.Split(toParse, "\t")
	if len(splitted) != 2 {
		return parsedEntry, errors.New("Invalid string (no single \\t)")
	}

	parsedEntry.Scope = splitted[0]

	splitted = strings.Split(splitted[1], "=")
	if len(splitted) != 2 {
		return parsedEntry, errors.New("Invalid string (no single =)")
	}

	parsedEntry.Name = splitted[0]
	parsedEntry.Value = splitted[1]

	return parsedEntry, nil
}

func (this *GitConfigMap) AddEntry(entry GitConfigEntry) {
	_, alreadyExists := (*this)[entry.Name]

	if alreadyExists {
		(*this)[entry.Name].Values[entry.Scope] = entry.Value
	} else {
		splitSectionKey := strings.SplitN(entry.Name, ".", 2)
		section := splitSectionKey[0]
		key := splitSectionKey[1]

		valuesMap := make(GitConfigPropValues)
		valuesMap[entry.Scope] = entry.Value

		(*this)[entry.Name] = GitConfigProp{
			Section: section,
			Key:     key,
			Type:    0,
			Values:  valuesMap,
		}
	}
}

func (this *GitConfigMap) AddEntries(entries []GitConfigEntry) {
	for _, entry := range entries {
		this.AddEntry(entry)
	}
}

func (this *GitConfigMap) AddLabel(label string) {
	_, alreadyExists := (*this)[label]

	if !alreadyExists {
		splitted := strings.SplitN(label, ".", 2)

		if len(splitted) != 2 {
			return
		}

		section := splitted[0]
		key := splitted[1]
		lowerLabel := strings.ToLower(label)

		(*this)[lowerLabel] = GitConfigProp{
			Section: section,
			Key:     strings.ToLower(key),
			Values:  make(map[string]string),
			Type:    0,
		}
	}
}

func (this *GitConfigMap) AddLabels(labels []string) {
	for _, label := range labels {
		this.AddLabel(label)
	}
}

func (this *GitConfigMap) ToSlice() []GitConfigProp {
	configSlice := slices.Collect(maps.Values((*this)))
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
