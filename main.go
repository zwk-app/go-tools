package main

import (
	"go-tools/tests"
	"time"
)

func main() {
	tests.RunTimerTests()
	time.Sleep(120 * time.Second)
}
