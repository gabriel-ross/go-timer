package main

import "github.com/gabriel-ross/timer-go"

func main() {
	cnf := timer.Config{
		MaxIntervals:        99,
		MaxIntervalDuration: timer.Interval{},
	}

	a := timer.New(cnf)
	a.Start()
}
