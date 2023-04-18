package internal

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var PRECISION = 300 * time.Millisecond

func TestRepeatTimer(t *testing.T) {
	cnf := Config{
		Intervals:       2,
		IntervalMinutes: 0,
		IntervalSeconds: 5,
		RestEnabled:     true,
		RestMinutes:     0,
		RestSeconds:     2,
		RestBeforeStart: true,
	}
	timer := NewRepeatCountdownTimer(cnf)
	intervalOrder := []string{}
	done := make(chan bool)

	expectedTime := 14 * time.Second // intervals*IntervalSeconds + Intervals*RestSeconds = 26

	go func() {
		fmt.Printf("test in progress expected time: %v", expectedTime)
		for {
			select {
			case name := <-timer.IntervalName():
				intervalOrder = append(intervalOrder, name)
			case <-done:
				fmt.Println()
				return
			case <-time.Tick(time.Second):
				fmt.Print(".")
			}
		}
	}()

	startTime := time.Now()
	timer.Start()
	done <- true
	actualTime := time.Since(startTime)
	diff := (expectedTime - actualTime).Abs()

	assert.LessOrEqualf(t, diff, PRECISION, "actual time %v is not within %v+-%v", actualTime, expectedTime, PRECISION)
	assert.ElementsMatch(t, intervalOrder, []string{"Starting", "Rest", "Interval 1/2", "Rest", "Interval 2/2"})
}
