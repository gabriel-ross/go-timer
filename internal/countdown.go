package internal

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Intervals       int
	IntervalMinutes int64
	IntervalSeconds int64
	RestEnabled     bool
	RestMinutes     int64
	RestSeconds     int64
	RestBeforeStart bool
}

type repeatTimer struct {
	cnf            Config
	shouldRest     bool
	shouldCancel   bool
	intervalNameC  chan string
	timeRemainingC chan string
	*countdownTimer
}

func NewRepeatCountdownTimer(cnf Config) *repeatTimer {
	// Format interval and rest seconds so that they don't exceed 59. Add
	// overflow to minutes
	for cnf.IntervalSeconds > 59 {
		cnf.IntervalMinutes++
		cnf.IntervalSeconds -= 59
	}
	for cnf.RestSeconds > 59 {
		cnf.RestMinutes++
		cnf.RestSeconds -= 59
	}

	return &repeatTimer{
		cnf:            cnf,
		shouldRest:     cnf.RestBeforeStart,
		shouldCancel:   false,
		intervalNameC:  make(chan string, 100),
		timeRemainingC: make(chan string, 100),
		countdownTimer: newCountdownTimer(),
	}
}

func (t *repeatTimer) Start() {
	t.reset()
	writeOut(t.intervalNameC, "Starting")

	if t.cnf.RestBeforeStart {
		writeOut(t.intervalNameC, "Rest")
		t.countdownTimer.runInterval(t.timeRemainingC, t.cnf.RestMinutes, t.cnf.RestSeconds)
	}

	interval := 1
	for interval <= t.cnf.Intervals && !t.shouldCancel {
		if t.shouldRest {
			writeOut(t.intervalNameC, "Rest")
			t.countdownTimer.runInterval(t.timeRemainingC, t.cnf.RestMinutes, t.cnf.RestSeconds)
		} else {
			writeOut(t.intervalNameC, fmt.Sprintf("Interval %d/%d", interval, t.cnf.Intervals))
			t.countdownTimer.runInterval(t.timeRemainingC, t.cnf.IntervalMinutes, t.cnf.IntervalSeconds)
			interval++
		}
		t.shouldRest = !t.shouldRest
	}
}

// reset resets all repeatTimer flags and clears all channels.
func (t *repeatTimer) reset() {
	t.shouldCancel = false
	t.shouldRest = false
	t.intervalNameC = make(chan string, 100)
	t.timeRemainingC = make(chan string, 100)
	t.countdownTimer = newCountdownTimer()
}

func (t *repeatTimer) IntervalName() <-chan string {
	return t.intervalNameC
}

func (t *repeatTimer) TimeRemaining() <-chan string {
	return t.timeRemainingC
}

// Cancel cancels the timer and writes zero time remaining to the
// time remaining channel.
func (t *repeatTimer) Cancel() {
	t.shouldCancel = true
	t.countdownTimer.cancel()
	writeOut(t.timeRemainingC, "00:00")
}

// Skip skips the current interval.
func (t *repeatTimer) Skip() {
	t.countdownTimer.cancel()
}

// RestartInterval restarts the current interval.
func (t *repeatTimer) RestartInterval() {
	t.countdownTimer.restart()
}

type countdownTimer struct {
	cancelC  chan bool
	pauseC   chan bool
	resumeC  chan bool
	restartC chan bool
}

func newCountdownTimer() *countdownTimer {
	return &countdownTimer{
		cancelC:  make(chan bool),
		pauseC:   make(chan bool),
		resumeC:  make(chan bool),
		restartC: make(chan bool),
	}
}

func (c *countdownTimer) runInterval(remainingC chan<- string, mins, secs int64) {
	writeOut(remainingC, formatTimeRemaining(mins, secs))

	// Decrement time before first tick, otherwise countdown is one second
	// longer than intended
	var remainingMins, remainingSecs int64
	switch {
	case secs > 0:
		remainingMins, remainingSecs = mins, secs-1
	case secs == 0 && mins > 0:
		remainingMins, remainingSecs = mins-1, 59
	default: // if minutes and seconds are both zero
		return
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.pauseC:
			select {
			case <-c.resumeC:
			case <-c.cancelC:
				return
			}
		case <-c.cancelC:
			return
		case <-c.restartC:
			remainingMins = mins
			remainingSecs = secs
		case <-ticker.C:
			writeOut(remainingC, formatTimeRemaining(remainingMins, remainingSecs))
			switch {
			case remainingMins > 0 && remainingSecs == 0:
				remainingMins--
				remainingSecs = 59
			case remainingSecs > 0:
				remainingSecs--
			case remainingMins == 0 && remainingSecs == 0:
				return
			}
		}
	}
}

func formatTimeRemaining(mins, secs int64) string {
	var builder strings.Builder
	if mins < 10 {
		builder.WriteString("0")
	}
	builder.WriteString(fmt.Sprintf("%d:", mins))
	if secs < 10 {
		builder.WriteString("0")
	}
	builder.WriteString(fmt.Sprintf("%d", secs))
	return builder.String()
}

func (c *countdownTimer) cancel() {
	c.cancelC <- true
}

func (c *countdownTimer) Pause() {
	c.pauseC <- true
}

func (c *countdownTimer) Resume() {
	c.resumeC <- true
}

func (c *countdownTimer) restart() {
	c.restartC <- true
}

// writeOut is a non-blocking helper function for writing outputs to channels.
// It attempts to write to the specified channel, but skips the write if the
// channel is full
func writeOut(ch chan<- string, out string) {
	select {
	case ch <- out:
	default:
	}
}
