package internal

import (
	"fmt"
	"strings"
	"time"
)

// TODO:
// add rest
// add resting before start

type Config struct {
	Intervals       int
	IntervalMinutes int
	IntervalSeconds int
	RestEnabled     bool
	RestMinutes     int
	RestSeconds     int
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

	if t.cnf.RestBeforeStart {
		t.intervalNameC <- "Starting"
		t.countdownTimer.runInterval(t.timeRemainingC, t.cnf.RestMinutes, t.cnf.RestSeconds)
	}

	interval := 1
	for interval <= t.cnf.Intervals && !t.shouldCancel {
		if t.shouldRest {
			t.intervalNameC <- "Rest"
			t.countdownTimer.runInterval(t.timeRemainingC, t.cnf.RestMinutes, t.cnf.RestSeconds)
		} else {
			t.intervalNameC <- fmt.Sprintf("Interval %d/%d", interval, t.cnf.Intervals)
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
	t.timeRemainingC <- "00:00"

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

func (c *countdownTimer) runInterval(remainingC chan<- string, mins, secs int) {
	remainingMins, remainingSecs := mins, secs
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
			remainingC <- formatTimeRemaining(remainingMins, remainingSecs)
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

func formatTimeRemaining(mins, secs int) string {
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
