package git

import "strings"

func FilterGitConfigProps(allProps []GitConfigProp, searchPhrase string) []GitConfigProp {
	if len(searchPhrase) == 0 {
		return allProps
	}

	result := make([]GitConfigProp, 0)

	for _, prop := range allProps {
		if strings.Contains(prop.String(), searchPhrase) {
			result = append(result, prop)
		}
	}

	return result
}
