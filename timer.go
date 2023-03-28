package timer

import (
	"fmt"
	"strings"
	"time"
)

type Timer struct {
	Intervals       int
	Rest            Interval
	IntervalLength  Interval
	RestSound       string
	StartSound      string
	RestBeforeStart bool
}

type Interval struct {
	Minutes int64
	Seconds int64
}

func (i Interval) Countdown() {
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

// TODO: Display remaining time and count down

func (t Timer) Start() {
	if t.RestBeforeStart {
		print("resting before start")
	}
}
