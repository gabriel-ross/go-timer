package timer

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Config struct {
	Intervals       int
	IntervalLength  Interval
	Rest            Interval
	SoundPath       string
	RestBeforeStart bool
}

type timer struct {
	intervals       int
	intervalLength  Interval
	rest            Interval
	sound           Stream
	restBeforeStart bool
	cancel          <-chan bool
	restart         <-chan bool
	skip            <-chan bool
}

type Stream struct {
	Streamer      beep.StreamSeekCloser
	StartPosition int
}

type Interval struct {
	Minutes int64
	Seconds int64
}

func New(cnf Config, cancel, restart, skip <-chan bool) (*timer, error) {
	var err error
	f, err := os.Open(cnf.SoundPath)
	if err != nil {
		return nil, err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return nil, err
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	return &timer{
		intervals:      cnf.Intervals,
		intervalLength: cnf.IntervalLength,
		rest:           cnf.Rest,
		sound: Stream{
			Streamer:      streamer,
			StartPosition: streamer.Position(),
		},
		restBeforeStart: cnf.RestBeforeStart,
		cancel:          cancel,
		restart:         restart,
		skip:            skip,
	}, nil
}

func (t timer) Close() {
	t.sound.Streamer.Close()
}

func (t timer) Start() string {
	if t.restBeforeStart {
		t.Countdown(t.rest, "Starting")
		t.PlaySound(t.sound, nil)
	}

	for i := 1; i <= t.intervals-1; i++ {
		t.Countdown(t.intervalLength, fmt.Sprintf("Interval %d/%d", i, t.intervals))
		t.PlaySound(t.sound, nil)
		t.Countdown(t.rest, "Rest")
		t.PlaySound(t.sound, nil)
	}
	t.Countdown(t.intervalLength, fmt.Sprintf("Interval %d/%d", t.intervals, t.intervals))
	t.PlaySound(t.sound, nil)
	return "timer done!"
}

func (t timer) PlaySound(stream Stream, done chan<- bool) {
	speaker.Clear()
	stream.Streamer.Seek(stream.StartPosition)
	speaker.Play(beep.Seq(stream.Streamer, beep.Callback(func() {
		if done != nil {
			done <- true
		}
	})))
}

func (t timer) Countdown(i Interval, name string) {
	fmt.Println(name)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	intervalTimer := time.NewTimer(time.Duration(i.Minutes)*time.Minute + time.Duration(i.Seconds)*time.Second)
	remaining := Interval{
		Minutes: i.Minutes,
		Seconds: i.Seconds,
	}

	for {
		select {
		case <-intervalTimer.C:
			fmt.Println("timer expired")
			return
		case <-t.skip:
			fmt.Println("timer skipped")
			return
		case <-t.cancel:
			// ?: should I use timer.Stop() here?
			return
		case <-t.restart:
			remaining.Minutes = i.Minutes
			remaining.Seconds = i.Seconds
		case <-ticker.C:
			fmt.Println(remaining.String())
			switch {
			case remaining.Minutes > 0 && remaining.Seconds == 0:
				remaining.Minutes--
				remaining.Seconds = 59
			case remaining.Seconds > 0:
				remaining.Seconds--
			case remaining.Minutes == 0 && remaining.Seconds == 0:
				fmt.Println("timer expired with zero minutes and seconds")
				return
			default:
				fmt.Println("encountered error")
			}
		}
	}

	// done := make(chan bool)
	// go func() {
	// 	time.Sleep((time.Duration(i.Minutes) * time.Minute) + (time.Duration(i.Seconds) * time.Second))
	// 	done <- true
	// }()

	// for remaining.Minutes >= 0 {
	// 	fmt.Println(remaining.String())
	// 	select {
	// 	case <-done:
	// 		print("00:00 done!\n")
	// 		return
	// 	case <-ticker.C:
	// 		switch {
	// 		case remaining.Seconds <= 0:
	// 			remaining.Minutes--
	// 			remaining.Seconds = 59
	// 		default:
	// 			remaining.Seconds--
	// 		}
	// 	}
	// }
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
