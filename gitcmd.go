package main

import (
	"log"
	"strings"
)

const (
	TYPE_STRING = "string"
)

type GitConfigEntry struct {
	EntryStr string
	Section  string
	Key      string
	Type     string
	Value    map[string]string
}

func GetGitConfigByLevel(source string) []GitConfigEntry {
	cmdOutputStr, err := executeCommand("git", "config", "--list", "--"+source)

	if err != nil {
		log.Println("Failed", err)
		return nil
	}

	entries := strings.Split(cmdOutputStr, "\n")

	result := make([]GitConfigEntry, len(entries))

	for index, entryStr := range entries {
		split := strings.SplitN(entryStr, "=", 2)
		location := strings.SplitN(split[0], ".", 2)

		result[index].EntryStr = entryStr
		result[index].Section = location[0]
		result[index].Key = split[0]
		result[index].Type = TYPE_STRING
		result[index].Value = make(map[string]string)
		result[index].Value[source] = split[1]
	}

	return result
}

func GetGitVersion() string {
	result, err := executeCommand("git", "version")
	if err != nil {
		return "git version ???"
	}

	return result
}
