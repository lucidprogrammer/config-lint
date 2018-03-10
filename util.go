package main

import (
	"path/filepath"
)

func unquoted(s string) string {
	if s[0] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func isAbsent(s string) bool {
	if s == "" || s == "null" || s == "[]" {
		return true
	}
	return false
}

func isPresent(s string) bool {
	return !isAbsent(s)
}

func listsIntersect(list1 []string, list2 []string) bool {
	for _, a := range list1 {
		for _, b := range list2 {
			if a == b {
				return true
			}
		}
	}
	return false
}

func shouldIncludeFile(patterns []string, filename string) bool {
	for _, pattern := range patterns {
		_, file := filepath.Split(filename)
		matched, err := filepath.Match(pattern, file)
		if err != nil {
			panic(err)
		}
		if matched {
			return true
		}
	}
	return false
}