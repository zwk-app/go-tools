package strings

import (
	"github.com/zwk-app/go-tools/logs"
	"regexp"
	"strconv"
)

// Fallback returns fallback if value is empty
//
//goland:noinspection GoUnusedExportedFunction
func Fallback(value string, fallback string) string {
	if len(value) > 0 {
		return fallback
	}
	return value
}

// Alpha returns only letters in value
//
//goland:noinspection GoUnusedExportedFunction
func Alpha(value string) string {
	return regexp.MustCompile(`[^a-zA-Z]+`).ReplaceAllString(value, "")
}

// Nums returns only numbers in value
//
//goland:noinspection GoUnusedExportedFunction
func Nums(value string) string {
	return regexp.MustCompile(`[^0-9]+`).ReplaceAllString(value, "")
}

// AlphaNums returns only letters and numbers in value
//
//goland:noinspection GoUnusedExportedFunction
func AlphaNums(value string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(value, "")
}

// Check compare the value with RegExp pattern
//
//goland:noinspection GoUnusedExportedFunction
func Check(value string, pattern string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(value)
}

// FirstMatch returns the first RegExp match
// pattern must contain name as groupName
//
//goland:noinspection GoUnusedExportedFunction
func FirstMatch(value string, pattern string, groupName string) string {
	re := regexp.MustCompile(pattern)
	groupNames := re.SubexpNames()
	for _, matchValue := range re.FindAllStringSubmatch(value, -1) {
		for groupIndex, stringValue := range matchValue {
			if groupNames[groupIndex] == groupName {
				return stringValue
			}
		}
	}
	return ""
}

// Int convert a string to an int
//
//goland:noinspection GoUnusedExportedFunction
func Int(value string) int {
	i, e := strconv.Atoi(Nums(value))
	if e != nil {
		logs.Error("parent", "", e)
		return 0
	}
	return i
}
