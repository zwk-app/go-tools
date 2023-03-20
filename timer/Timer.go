package timer

import (
	"fmt"
	"go-tools/logs"
	"sort"
	"time"
)

const defaultDelayBeforeNext int64 = -15

//goland:noinspection GoNameStartsWithPackageName
type Timer struct {
	Targets []TargetInfo
	Next    *TargetInfo
	Current struct {
		//Time chan time.Time
		Time time.Time
		//Text chan string
		Text string
	}
	Remaining struct {
		Duration chan time.Duration
		//Text     chan string
		Text string
		//Seconds  chan int64
		Seconds int64
	}
	Alarm struct {
		Callback func(name string, alarm string)
	}
	Alerts struct {
		Callback func(name string, remaining int64)
	}
	running bool
}

type TargetInfo struct {
	Time struct {
		Object time.Time
		String TimeString
		Text   string
	}
	Name     string
	Alarm    string
	OnlyOnce bool
}

func (r *TargetInfo) String() string {
	return fmt.Sprintf("%-16s %-10s (%s) <%t>\n", r.Name, r.Time.Text, r.Alarm, r.OnlyOnce)
}

var timer *Timer = nil

//goland:noinspection GoNameStartsWithPackageName,GoUnusedExportedFunction
func getTimer() *Timer {
	if timer == nil {
		timer = new(Timer)
		timer.Alerts.Callback = nil
		timer.Alarm.Callback = nil
		timer.running = false
		//timer.Current.Time = make(chan time.Time)
		//timer.Current.Text = make(chan string)
		currentTime := time.Now()
		timer.Current.Time = currentTime
		timer.Current.Text = TimeTextFromObject(currentTime)

		timer.Remaining.Duration = make(chan time.Duration)
		//timer.Remaining.Text = make(chan string)
		//timer.Remaining.Seconds = make(chan int64)
		timer.timerLoop()
	}
	return timer
}

func (r *Timer) isRunning() bool {
	if r.Next != nil && r.running == false {
		Start()
	}
	if r.Next == nil && r.running == true {
		Stop()
	}
	return r.running
}

func (r *Timer) timerLoop() {
	go func() {
		time.Sleep(800 * time.Millisecond)
		logs.Debug("Timer->Loop", "Start", nil)
		for {
			currentTime := time.Now()
			r.Current.Time = currentTime
			r.Current.Text = TimeTextFromObject(currentTime)
			if r.isRunning() {
				duration := r.Next.Time.Object.Sub(currentTime).Round(time.Second)
				r.Remaining.Duration <- duration
				r.Remaining.Text = DelayTextFromObject(duration)
				r.Remaining.Seconds = int64(duration / time.Second)
			}
			time.Sleep(250 * time.Millisecond)
		}
		logs.Debug("Timer->Loop", "Stop", nil)
	}()
}

func (r *Timer) alertLoop() {
	var currentCheck int64 = 0
	var lastCheck int64 = 0
	var lastCheckDiff int64 = 0
	go func() {
		logs.Debug("Timer->AlertLoop", "Start", nil)
		for remaining := range r.Remaining.Duration {
			currentCheck = DelaySecondsFromObject(remaining.Round(time.Second))
			lastCheckDiff = lastCheck - currentCheck
			if lastCheckDiff > 1 && lastCheckDiff <= 10 {
				// some lag could happen?
				currentCheck = lastCheck + 1
			}
			if currentCheck < lastCheck {
				// only once per second
				go r.alertCheck(&remaining)
			}
			lastCheck = currentCheck
			time.Sleep(200 * time.Millisecond)
		}
		logs.Debug("Timer->AlertLoop", "Stop", nil)
	}()
}

func (r *Timer) alertCheck(remaining *time.Duration) {
	seconds := DelaySecondsFromObject(*remaining)
	if seconds < 60 {
		switch seconds {
		case 0:
			r.alarm()
		case 2, 4, 6, 8:
			r.alert(seconds)
		case 10, 20, 30:
			r.alert(seconds)
		}
		if seconds < defaultDelayBeforeNext {
			if r.Next.OnlyOnce {
				r.delTarget(r.Next)
			}
			r.nextTarget()
		}
	} else if seconds < 310 && seconds%60 == 0 {
		// every 1m (60) if T <= 5m (300)
		r.alert(seconds)
	} else if seconds < 910 && seconds%300 == 0 {
		// every 5m (300) if T <= 15m (900)
		r.alert(seconds)
	} else if seconds < 1810 && seconds%600 == 0 {
		// every 10m (600) if T <= 30m (1800)
		r.alert(seconds)
	} else if seconds < 10810 && seconds%3600 == 0 {
		// every 1h (3600) if T <= 3h (10800)
		r.alert(seconds)
	}
}

func (r *Timer) alert(seconds int64) {
	if r.Alerts.Callback != nil {
		go r.Alerts.Callback(r.Next.Name, seconds)
	}
}

func (r *Timer) alarm() {
	if r.Alarm.Callback != nil {
		go r.Alarm.Callback(r.Next.Name, r.Next.Alarm)
	}
}

func (r *Timer) setNextTarget(index int) {
	if index >= 0 && index <= len(r.Targets) {
		r.Next = &r.Targets[index]
		r.Next.Time.Object, _ = r.Next.Time.String.NextObject()
		r.Next.Time.Text = r.Next.Time.String.Text()
	}
}

func (r *Timer) nextTarget() {
	r.Next = nil
	current := TimeStringFromObject(time.Now())
	if len(r.Targets) > 0 {
		sort.Slice(r.Targets, func(i, j int) bool { return r.Targets[i].Time.String < r.Targets[j].Time.String })
		for i, v := range r.Targets {
			if current < v.Time.String {
				r.setNextTarget(i)
				break
			}
		}
		if r.Next == nil {
			r.setNextTarget(0)
		}
		logs.Debug("Timer->NextTarget", r.Next.Time.Text, nil)
	}
}

func (r *Timer) getTargetIndex(v *TargetInfo) int {
	for i, t := range r.Targets {
		if t.Time.String == v.Time.String {
			return i
		}
	}
	return -1
}

func (r *Timer) delTarget(v *TargetInfo) {
	logs.Debug("Timer->DelTarget", v.String(), nil)
	if i := r.getTargetIndex(v); i >= 0 {
		r.Targets = append(r.Targets[:i], r.Targets[i+1:]...)
	}
}

func (r *Timer) addTarget(v *TargetInfo) {
	logs.Debug("Timer->AddTarget", v.String(), nil)
	if i := r.getTargetIndex(v); i >= 0 {
		r.Targets[i].Name = v.Name
		r.Targets[i].Alarm = v.Alarm
	} else {
		r.Targets = append(r.Targets, *v)
	}
}

//goland:noinspection GoUnusedExportedFunction
func SetAlarmCallback(callback func(name string, alarm string)) {
	getTimer().Alarm.Callback = callback
}

//goland:noinspection GoUnusedExportedFunction
func SetAlertCallback(callback func(name string, remaining int64)) {
	getTimer().Alerts.Callback = callback
}

//goland:noinspection GoUnusedExportedFunction
func AddTargetTime(targetTime TimeString, name string, alarm string) error {
	logs.Debug("Timer->AddTargetTime", fmt.Sprintf("%-16s %-10s (%s)", name, targetTime, alarm), nil)
	if targetTime.Validate() {
		target := new(TargetInfo)
		target.Time.Object = time.Time{}
		target.Time.String = targetTime
		target.Time.Text = targetTime.Text()
		target.Name = name
		target.Alarm = alarm
		target.OnlyOnce = false
		getTimer().addTarget(target)
		return nil
	}
	return fmt.Errorf("invalid time string '%s'", targetTime)
}

//goland:noinspection GoUnusedExportedFunction
func AddTargetDelay(targetDelay DelayString, name string, alarm string) error {
	logs.Debug("Timer->AddTargetDelay", fmt.Sprintf("%-16s %-10s (%s)", name, targetDelay, alarm), nil)
	if targetDelay.Validate() {
		v := TimeStringFromObject(time.Now().Add(targetDelay.DelayObject()))
		target := new(TargetInfo)
		target.Time.Object = time.Time{}
		target.Time.String = v
		target.Time.Text = v.Text()
		target.Name = name
		target.Alarm = alarm
		target.OnlyOnce = true
		getTimer().addTarget(target)
		return nil
	}
	return fmt.Errorf("invalid time string '%s'", targetDelay)
}

//goland:noinspection GoUnusedExportedFunction
func Start() {
	logs.Debug("Timer->Start", "", nil)
	r := getTimer()
	if r.running == false && len(r.Targets) > 0 {
		r.nextTarget()
		r.running = true
		r.alertLoop()
	}
}

//goland:noinspection GoUnusedExportedFunction
func Stop() {
	logs.Debug("Timer->Stop", "", nil)
	r := getTimer()
	r.Next = nil
	r.running = false
}
