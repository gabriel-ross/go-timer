package main

import "github.com/gabriel-ross/timer-go"

func main() {
	a := timer.New(timer.Config{
		MaxIntervals: 99,
		MaxTimerMins: 99,
		MaxTimerSecs: 59,
	})
	a.Start()
}
