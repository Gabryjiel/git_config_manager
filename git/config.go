package git

import (
	"fmt"
)

const (
	GIT_CONFIG_PROP__TYPE_STRING = iota
)

type GitConfigPropValues map[string]string

type GitConfigProp struct {
	Section string
	Key     string
	Type    int
	Values  GitConfigPropValues
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
