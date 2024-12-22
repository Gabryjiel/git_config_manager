package git

import (
	"fmt"
)

const (
	GIT_CONFIG_PROP__TYPE_STRING = iota
)

type GitConfigProp struct {
	Section string
	Key     string
	Type    int
	Values  GitConfigEntryValues
}

func (this *GitConfigProp) String() string {
	return fmt.Sprintf("%s.%s=%s / %s / %s", this.Section, this.Key, this.Values.System, this.Values.Global, this.Values.Local)
}

func (this *GitConfigProp) GetName() string {
	return fmt.Sprintf("%s.%s", this.Section, this.Key)
}

type GitConfigEntryValues struct {
	Local  string
	Global string
	System string
}

func (this *GitConfigEntryValues) append(values GitConfigEntryValues) {
	if len(values.Local) != 0 {
		this.Local = values.Local
	}
	if len(values.Global) != 0 {
		this.Global = values.Global
	}
	if len(values.System) != 0 {
		this.System = values.System
	}
}

func CreateConfigMap(props ...[]GitConfigProp) []GitConfigProp {
	gitConfig := []GitConfigProp{
		{Section: "core", Key: "bare", Type: GIT_CONFIG_PROP__TYPE_STRING},
		{Section: "core", Key: "editor", Type: GIT_CONFIG_PROP__TYPE_STRING},
		{Section: "core", Key: "filemode", Type: GIT_CONFIG_PROP__TYPE_STRING},
		{Section: "user", Key: "email", Type: GIT_CONFIG_PROP__TYPE_STRING},
		{Section: "user", Key: "name", Type: GIT_CONFIG_PROP__TYPE_STRING},
	}

	for _, entryGroup := range props {
		for _, entry := range entryGroup {
			for index, existing := range gitConfig {
				if existing.Section == entry.Section && existing.Key == entry.Key {
					gitConfig[index].Values.append(entry.Values)
					break
				}

				if index == len(gitConfig)-1 {
					gitConfig = append(gitConfig, entry)
				}
			}
		}
	}

	return gitConfig
}
