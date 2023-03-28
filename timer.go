package timer

import (
	"fmt"
	"strings"
	"time"
)

type Timer struct {
	Intervals       int
	IntervalLength  Interval
	Rest            Interval
	IntervalSound   string
	RestSound       string
	RestBeforeStart bool
}

type Interval struct {
	Minutes int64
	Seconds int64
}

// TODO: Display remaining time and count down

func (t Timer) Start() {
	if t.RestBeforeStart {
		print("resting before start")
	}

	for i := 1; i <= t.Intervals-1; i++ {
		t.Countdown(t.IntervalLength, fmt.Sprintf("Interval %d/%d", i, t.Intervals), t.IntervalSound)
		t.Countdown(t.Rest, "Rest", t.RestSound)
	}
	t.Countdown(t.IntervalLength, fmt.Sprintf("Interval %d/%d", t.Intervals, t.Intervals), t.IntervalSound)
}

// TODO: add end sound
func (t Timer) Countdown(i Interval, name, endSound string) {
	fmt.Println(name)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	done := make(chan bool)

	go func() {
		time.Sleep((time.Duration(i.Minutes) * time.Minute) + (time.Duration(i.Seconds) * time.Second))
		done <- true
	}()

	for i.Minutes >= 0 {
		select {
		case <-done:
			print("00:00 done!\n")
			return
		case <-ticker.C:
			switch {
			case i.Seconds <= 0:
				i.Minutes--
				i.Seconds = 59
			default:
				i.Seconds--
			}
			fmt.Println(i.String())
		}
	}
}

func (i *Interval) String() string {
	var builder strings.Builder
	if i.Minutes < 10 {

		builder.WriteString("0")
	}
	builder.WriteString(fmt.Sprintf("%d:", i.Minutes))

	if i.Seconds < 10 {

		builder.WriteString("0")
	}
	builder.WriteString(fmt.Sprintf("%d", i.Seconds))

	return builder.String()
}
