package timer

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type TimeString string

func TimeStringNums(v string) string {
	return regexp.MustCompile(`[^0-9]+`).ReplaceAllString(v, "")
}

func TimeValidate(h int, m int, s int) bool {
	return h >= 0 && h <= 23 && m >= 0 && m <= 59 && s >= 0 && s <= 59
}

//goland:noinspection GoUnusedExportedFunction,GoRedundantConversion
func TimeFromObject(v time.Time) (int, int, int) {
	return int(v.Hour()), int(v.Minute()), int(v.Second())
}

func TimeFromString(v string) (int, int, int, error) {
	timeString := TimeStringNums(v)
	hasError := false
	var he, me, se error
	h := 0
	m := 0
	s := 0
	switch len(timeString) {
	case 1, 2:
		h, he = strconv.Atoi(timeString)
		if he != nil {
			hasError = true
		}
	case 3:
		h, he = strconv.Atoi(timeString[0:1])
		m, me = strconv.Atoi(timeString[1:3])
		if he != nil || me != nil {
			hasError = true
		}
	case 4:
		h, he = strconv.Atoi(timeString[0:2])
		m, me = strconv.Atoi(timeString[2:4])
		if he != nil || me != nil {
			hasError = true
		}
	case 5:
		h, he = strconv.Atoi(timeString[0:1])
		m, me = strconv.Atoi(timeString[1:3])
		s, se = strconv.Atoi(timeString[3:5])
		if he != nil || me != nil || se != nil {
			hasError = true
		}
	case 6:
		h, he = strconv.Atoi(timeString[0:2])
		m, me = strconv.Atoi(timeString[2:4])
		s, se = strconv.Atoi(timeString[4:6])
		if he != nil || me != nil || se != nil {
			hasError = true
		}
	}
	if !hasError {
		if TimeValidate(h, m, s) {
			return h, m, s, nil
		}
	}
	return 0, 0, 0, fmt.Errorf("invalid time string '%s'", timeString)
}

func NextObjectFromTime(h int, m int, s int) (time.Time, error) {
	if TimeValidate(h, m, s) {
		currentTime := time.Now()
		tmpTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), h, m, s, 0, time.Local)
		if tmpTime.Before(currentTime) {
			tmpTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, h, m, s, 0, time.Local)
		}
		return tmpTime, nil
	}
	return time.Time{}, fmt.Errorf("invalid time '%02d:%02d:%02d'", h, m, s)
}

//goland:noinspection GoUnusedExportedFunction
func NextObjectFromTimeString(v string) (time.Time, error) {
	if h, m, s, e := TimeFromString(v); e == nil {
		return NextObjectFromTime(h, m, s)
	}
	return time.Time{}, fmt.Errorf("invalid time string '%s'", v)
}

//goland:noinspection GoUnusedExportedFunction
func TimeSecondsFromObject(v time.Time) int64 {
	return v.Unix()
}

func TimeSecondsFromTime(h int, m int, s int) (int64, error) {
	if tmpTime, e := NextObjectFromTime(h, m, s); e == nil {
		return tmpTime.Unix(), nil
	}
	return 0, fmt.Errorf("invalid time '%02d:%02d:%02d'", h, m, s)
}

//goland:noinspection GoUnusedExportedFunction
func TimeSecondsFromString(v string) (int64, error) {
	if h, m, s, e := TimeFromString(v); e == nil {
		if seconds, e := TimeSecondsFromTime(h, m, s); e == nil {
			return seconds, nil
		}
	}
	return 0, fmt.Errorf("invalid time string '%s'", v)
}

//goland:noinspection GoUnusedExportedFunction
func TimeStringFromObject(v time.Time) TimeString {
	h, m, s := v.Hour(), v.Minute(), v.Second()
	timeString, _ := TimeStringFromTime(h, m, s)
	return timeString
}

func TimeStringFromTime(h int, m int, s int) (TimeString, error) {
	if TimeValidate(h, m, s) {
		return TimeString(fmt.Sprintf("%02d:%02d:%02d", h, m, s)), nil
	}
	//goland:noinspection GoRedundantConversion
	return TimeString(""), fmt.Errorf("invalid time '%02d:%02d:%02d'", h, m, s)
}

//goland:noinspection GoUnusedExportedFunction
func TimeStringFromString(v string) (TimeString, error) {
	if h, m, s, e := TimeFromString(v); e == nil {
		return TimeStringFromTime(h, m, s)
	}
	//goland:noinspection GoRedundantConversion
	return TimeString(""), fmt.Errorf("invalid time string '%s'", v)
}

func (r TimeString) Validate() bool {
	if h, m, s, e := r.Time(); e == nil {
		return TimeValidate(h, m, s)
	}
	return false
}

func (r TimeString) Time() (int, int, int, error) {
	if h, m, s, e := TimeFromString(string(r)); e == nil {
		return h, m, s, nil
	}
	return 0, 0, 0, fmt.Errorf("invalid time string '%s'", string(r))
}

func (r TimeString) Text() string {
	return TimeTextFromString(string(r))
}

func (r TimeString) NextObject() (time.Time, error) {
	return NextObjectFromTimeString(string(r))
}

func TimeTextFromTime(h int, m int, s int) string {
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

//goland:noinspection GoUnusedExportedFunction
func TimeTextFromObject(v time.Time) string {
	return TimeTextFromTime(v.Hour(), v.Minute(), v.Second())
}

//goland:noinspection GoUnusedExportedFunction
func TimeTextFromString(v string) string {
	if h, m, s, e := TimeFromString(v); e == nil {
		return TimeTextFromTime(h, m, s)
	}
	return fmt.Sprintf("--:--:--")
}
