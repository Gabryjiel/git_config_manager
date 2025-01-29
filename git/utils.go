package git

import "strings"

func FilterGitConfigProps(allProps []GitConfigProp, searchPhrase string, onlyWithValue bool) []GitConfigProp {
	result := make([]GitConfigProp, 0)

	for _, prop := range allProps {
		if strings.Contains(prop.String(), searchPhrase) {
			if onlyWithValue && len(prop.Values) == 0 {
				continue
			}

			result = append(result, prop)
		}
	}

	return result
}
