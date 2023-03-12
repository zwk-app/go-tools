package timer

import (
	"fmt"
	"strconv"
	"time"
)

type DelayString string

func DelayValidate(h int, m int, s int) bool {
	return h >= 0 && m >= 0 && m <= 59 && s >= 0 && s <= 59
}

func DelayFromObject(v time.Duration) (int, int, int) {
	d := v.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return int(h), int(m), int(s)
}

func DelayFromString(v string) (int, int, int, error) {
	delayString := TimeStringNums(v)
	hasError := false
	var he, me, se error
	h := 0
	m := 0
	s := 0
	switch len(delayString) {
	case 1, 2:
		s, se = strconv.Atoi(delayString)
		if se != nil {
			hasError = true
		}
	case 3:
		m, me = strconv.Atoi(delayString[0:1])
		s, se = strconv.Atoi(delayString[1:3])
		if me != nil || se != nil {
			hasError = true
		}
	case 4:
		m, me = strconv.Atoi(delayString[0:2])
		s, se = strconv.Atoi(delayString[2:4])
		if me != nil || se != nil {
			hasError = true
		}
	case 5:
		h, he = strconv.Atoi(delayString[0:1])
		m, me = strconv.Atoi(delayString[1:3])
		s, se = strconv.Atoi(delayString[3:5])
		if he != nil || me != nil || se != nil {
			hasError = true
		}
	case 6:
		h, he = strconv.Atoi(delayString[0:2])
		m, me = strconv.Atoi(delayString[2:4])
		s, se = strconv.Atoi(delayString[4:6])
		if he != nil || me != nil || se != nil {
			hasError = true
		}
	}
	if !hasError {
		if DelayValidate(h, m, s) {
			return h, m, s, nil
		}
	}
	return 0, 0, 0, fmt.Errorf("invalid delay string '%s'", delayString)
}

//goland:noinspection GoUnusedExportedFunction
func DelaySecondsFromObject(v time.Duration) int64 {
	return int64(v.Round(time.Second) / time.Second)
}

func DelaySecondsFromDelay(h int, m int, s int) (int64, error) {
	if DelayValidate(h, m, s) {
		return int64(s + (m * 60) + (h * 3600)), nil
	}
	return 0, fmt.Errorf("invalid delay '%02d:%02d:%02d'", h, m, s)
}

//goland:noinspection GoUnusedExportedFunction
func DelaySecondsFromString(v string) (int64, error) {
	if h, m, s, e := DelayFromString(v); e == nil {
		if seconds, e := DelaySecondsFromDelay(h, m, s); e == nil {
			return seconds, nil
		}
	}
	return 0, fmt.Errorf("invalid delay string '%s'", v)
}

//goland:noinspection GoUnusedExportedFunction
func DelayStringFromObject(v time.Duration) DelayString {
	h, m, s := DelayFromObject(v)
	delayString, _ := DelayStringFromDelay(h, m, s)
	return delayString
}

func DelayStringFromDelay(h int, m int, s int) (DelayString, error) {
	if DelayValidate(h, m, s) {
		if h > 0 {
			return DelayString(fmt.Sprintf("%02d:%02d:%02d", h, m, s)), nil
		} else if m > 0 {
			return DelayString(fmt.Sprintf("%02d:%02d", m, s)), nil
		} else if s > 0 {
			return DelayString(fmt.Sprintf("%02d", s)), nil
		}
	}
	//goland:noinspection GoRedundantConversion
	return DelayString(""), fmt.Errorf("invalid delay '%02d:%02d:%02d'", h, m, s)
}

//goland:noinspection GoUnusedExportedFunction
func DelayStringFromString(v string) (DelayString, error) {
	if h, m, s, e := DelayFromString(v); e == nil {
		return DelayStringFromDelay(h, m, s)
	}
	//goland:noinspection GoRedundantConversion
	return DelayString(""), fmt.Errorf("invalid delay string '%s'", v)
}

func (r DelayString) Validate() bool {
	if h, m, s, e := r.Delay(); e == nil {
		return DelayValidate(h, m, s)
	}
	return false
}

func (r DelayString) Delay() (int, int, int, error) {
	if h, m, s, e := DelayFromString(string(r)); e == nil {
		return h, m, s, nil
	}
	return 0, 0, 0, fmt.Errorf("invalid delay string '%s'", string(r))
}

func DelayTextFromDelay(h int, m int, s int) string {
	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	} else if m > 0 {
		return fmt.Sprintf("%02d:%02d", m, s)
	} else if s > 0 {
		return fmt.Sprintf("%02d", s)
	} else if s > -30 {
		return fmt.Sprintf(">> %02d <<", s)
	}
	return fmt.Sprintf("--")
}

//goland:noinspection GoUnusedExportedFunction
func DelayTextFromObject(v time.Duration) string {
	h, m, s := DelayFromObject(v)
	return DelayTextFromDelay(h, m, s)
}

//goland:noinspection GoUnusedExportedFunction
func DelayTextFromString(v string) string {
	if h, m, s, e := DelayFromString(v); e == nil {
		return DelayTextFromDelay(h, m, s)
	}
	return fmt.Sprintf("--")
}
