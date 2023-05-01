package tests

import (
	"fmt"
	"time"
	"zwk-tools/logs"
	"zwk-tools/timer"
)

func TestTimerAlarmCallback(name string, alarm string) {
	logs.Debug("Tests->Timer->Alarm", fmt.Sprintf("%s: %s", name, alarm), nil)
}

func TestTimerAlertCallback(name string, remaining int64) {
	logs.Debug("Tests->Timer->Alert", fmt.Sprintf("%s: %08d", name, remaining), nil)
}
func TestTimerTargetTime() {
	_ = timer.AddTargetTime(
		timer.TimeStringFromObject(time.Now().Add(15*time.Second)),
		"TestTargetTime",
		"TaDaa!")
}

func TestTimerTargetDelay() {
	_ = timer.AddTargetDelay(
		"100",
		"TestTargetDelay",
		"TaDaa!")
}

func RunTimerTests() {
	logs.SetLevelDebug()
	timer.SetAlarmCallback(TestTimerAlarmCallback)
	timer.SetAlertCallback(TestTimerAlertCallback)
	TestTimerTargetTime()
	TestTimerTargetDelay()
	timer.Start()
	time.Sleep(120 * time.Second)
}
