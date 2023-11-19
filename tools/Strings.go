package tools

import (
	"errors"
	"fmt"
	"github.com/zwk-app/zwk-tools/logs"
	"regexp"
	"strconv"
)

// StringFallback returns fallback if value is empty
//
//goland:noinspection GoUnusedExportedFunction
func StringFallback(value string, fallback string) string {
	if len(value) > 0 {
		return fallback
	}
	return value
}

// StringAlpha returns only letters in value
//
//goland:noinspection GoUnusedExportedFunction
func StringAlpha(value string) string {
	return regexp.MustCompile(`[^a-zA-Z]+`).ReplaceAllString(value, "")
}

// StringNums returns only numbers in value
//
//goland:noinspection GoUnusedExportedFunction
func StringNums(value string) string {
	return regexp.MustCompile(`[^0-9]+`).ReplaceAllString(value, "")
}

// StringAlphaNums returns only letters and numbers in value
//
//goland:noinspection GoUnusedExportedFunction
func StringAlphaNums(value string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(value, "")
}

// StringCheck compare the value with RegExp pattern
//
//goland:noinspection GoUnusedExportedFunction
func StringCheck(value string, pattern string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(value)
}

// StringFirstMatch returns the first RegExp match
// pattern must contain name as groupName
//
//goland:noinspection GoUnusedExportedFunction
func StringFirstMatch(value string, pattern string, groupName string) string {
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

// StringToInt convert a string to an int
//
//goland:noinspection GoUnusedExportedFunction
func StringToInt(value string) int {
	i, e := strconv.Atoi(StringNums(value))
	if e != nil {
		logs.Error("parent", "", e)
		return 0
	}
	return i
}

// StringToBool convert a string to a boolean
//
//goland:noinspection GoUnusedExportedFunction
func StringToBool(value string) bool {
	switch value {
	case "true", "True", "TRUE":
		return true
	case "false", "False", "FALSE":
		return false
	default:
		logs.Error("parent", "", errors.New(fmt.Sprintf("StringToBool: invalid string value '%s'", value)))
		return false
	}
}
