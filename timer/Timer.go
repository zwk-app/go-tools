package timer

import (
	"fmt"
	"sort"
	"time"
)

const defaultDelayBeforeNext int64 = -20

//goland:noinspection GoNameStartsWithPackageName
type Timer struct {
	Targets []*TargetInfo
	Next    *TargetInfo
	Current struct {
		Time chan time.Time
		Text chan string
	}
	Remaining struct {
		Duration chan time.Duration
		Text     chan string
		Seconds  chan int64
	}
	Alerts struct {
		Callback func(*time.Duration)
	}
	Alarm struct {
		Callback func(*time.Duration)
	}
	running bool
}

type TargetInfo struct {
	Time struct {
		Object time.Time
		String TimeString
		Text   string
	}
	Name  string
	Alarm string
}

var timer *Timer = nil

//goland:noinspection GoNameStartsWithPackageName,GoUnusedExportedFunction
func getTimer() *Timer {
	if timer == nil {
		timer = new(Timer)
		timer.Alerts.Callback = nil
		timer.Alarm.Callback = nil
		timer.running = false
		timer.Current.Time = make(chan time.Time)
		timer.Current.Text = make(chan string)
		timer.Remaining.Duration = make(chan time.Duration)
		timer.Remaining.Text = make(chan string)
		timer.Remaining.Seconds = make(chan int64)
		timer.timerLoop()
	}
	return timer
}

func (r *Timer) isRunning() bool {
	if r.Next != nil && r.running == false {
		Start()
	} else if r.Next == nil && r.running == true {
		Stop()
	}
	return r.running
}

func (r *Timer) timerLoop() {
	go func() {
		time.Sleep(1500 * time.Millisecond)
		for {
			currentTime := time.Now()
			r.Current.Time <- currentTime
			r.Current.Text <- TimeTextFromObject(&currentTime)
			if r.isRunning() {
				duration := r.Next.Time.Object.Sub(currentTime).Round(time.Second)
				r.Remaining.Duration <- duration
				r.Remaining.Text <- DelayTextFromObject(duration)
				r.Remaining.Seconds <- int64(duration / time.Second)
			}
			time.Sleep(250 * time.Millisecond)
		}
	}()
}

func (r *Timer) alertLoop() {
	var currentCheck int64 = 0
	var lastCheck int64 = 0
	var lastCheckDiff int64 = 0
	go func() {
		for remaining := range r.Remaining.Duration {
			currentCheck = DelaySecondsFromObject(remaining)
			lastCheckDiff = lastCheck - currentCheck
			if 0 < lastCheckDiff && lastCheckDiff < 10 {
				// some lag could happen?
				currentCheck = lastCheck + 1
			}
			if lastCheckDiff != 0 {
				// only once per second
				r.alertCheck(&remaining)
			}
			lastCheck = currentCheck
			time.Sleep(250 * time.Millisecond)
		}
	}()
}

func (r *Timer) alertCheck(remaining *time.Duration) {
	seconds := DelaySecondsFromObject(*remaining)
	if seconds < 60 {
		switch seconds {
		case 0:
			r.alarm(remaining)
		case 2, 4, 6, 8:
			r.alert(remaining)
		case 10, 20, 30:
			r.alert(remaining)
		}
		if seconds < defaultDelayBeforeNext {
			r.nextTarget()
		}
	} else if seconds < 310 && seconds%60 == 0 {
		// every 1m (60) if T <= 5m (300)
		r.alert(remaining)
	} else if seconds < 910 && seconds%300 == 0 {
		// every 5m (300) if T <= 15m (900)
		r.alert(remaining)
	} else if seconds < 1810 && seconds%600 == 0 {
		// every 10m (600) if T <= 30m (1800)
		r.alert(remaining)
	} else if seconds < 10810 && seconds%3600 == 0 {
		// every 1h (3600) if T <= 3h (10800)
		r.alert(remaining)
	}
}

func (r *Timer) alert(remaining *time.Duration) {
	if r.Alerts.Callback != nil {
		r.Alerts.Callback(remaining)
	}
}

func (r *Timer) alarm(remaining *time.Duration) {
	if r.Alarm.Callback != nil {
		r.Alarm.Callback(remaining)
	}
}

func (r *Timer) nextTarget() {
	current := TimeStringFromObject(time.Now())
	if len(r.Targets) > 0 {
		sort.Slice(r.Targets, func(i, j int) bool { return r.Targets[i].Time.String < r.Targets[j].Time.String })
		for i, v := range r.Targets {
			if current < v.Time.String {
				r.Next = r.Targets[i]
				r.Next.Time.Object, _ = r.Next.Time.String.NextObject()
				r.Next.Time.Text = r.Next.Time.String.Text()
				return
			}
		}
		r.Next = r.Targets[0]
	} else {
		r.Next = nil
	}
}

func (r *Timer) addTarget(v *TargetInfo) {
	found := false
	for i, t := range r.Targets {
		if t.Time.String == v.Time.String {
			r.Targets[i].Name = v.Name
			r.Targets[i].Alarm = v.Alarm
			found = true
		}
	}
	if !found {
		r.Targets = append(r.Targets, v)
	}
}

//goland:noinspection GoUnusedExportedFunction
func SetTarget(targetTime TimeString, targetName string, targetAlarm string) error {
	if targetTime.Validate() {
		target := new(TargetInfo)
		target.Time.Object = time.Time{}
		target.Time.String = targetTime
		target.Time.Text = targetTime.Text()
		target.Name = targetName
		target.Alarm = targetAlarm
		getTimer().addTarget(target)
		return nil
	}
	return fmt.Errorf("invalid time string '%s'", targetTime)
}

//goland:noinspection GoUnusedExportedFunction
func Start() {
	r := getTimer()
	if r.running == false && len(r.Targets) > 0 {
		r.nextTarget()
	}
}

//goland:noinspection GoUnusedExportedFunction
func Stop() {
	r := getTimer()
	r.running = false
	r.Next = nil
}
