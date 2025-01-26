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
	Values  map[string]string
}

func (this *GitConfigProp) String() string {
	localValue, _ := this.Values["local"]
	globalValue, _ := this.Values["global"]
	systemValue, _ := this.Values["system"]

	return fmt.Sprintf("%s.%s=%s / %s / %s", this.Section, this.Key, systemValue, globalValue, localValue)
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
