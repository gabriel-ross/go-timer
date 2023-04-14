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

type timerConfig struct {
	Intervals       int
	IntervalLength  Interval
	Rest            Interval
	SoundPath       string
	RestBeforeStart bool
}

type countdownTimer struct {
	intervals       int
	intervalLength  Interval
	rest            Interval
	sound           Stream
	restBeforeStart bool
	pauseC          chan bool
	resumeC         chan bool
	skipC           chan bool
	cancelC         chan bool
	restartC        chan bool
	doneC           chan bool
	timeRemaining   chan string
	intervalName    chan string
	paused          bool
	cancelled       bool
}

type Stream struct {
	Streamer      beep.StreamSeekCloser
	StartPosition int
}

type Interval struct {
	Minutes int
	Seconds int
}

func NewTimer(cnf timerConfig) (*countdownTimer, error) {
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
	return &countdownTimer{
		intervals:      cnf.Intervals,
		intervalLength: cnf.IntervalLength,
		rest:           cnf.Rest,
		sound: Stream{
			Streamer:      streamer,
			StartPosition: streamer.Position(),
		},
		restBeforeStart: cnf.RestBeforeStart,
		pauseC:          make(chan bool),
		resumeC:         make(chan bool),
		cancelC:         make(chan bool),
		restartC:        make(chan bool),
		skipC:           make(chan bool),
		doneC:           make(chan bool),
		timeRemaining:   make(chan string, 61),
		intervalName:    make(chan string, cnf.Intervals),
		paused:          false,
		cancelled:       false,
	}, nil
}

func (t *countdownTimer) Close() {
	t.sound.Streamer.Close()
}

// Start starts the countdownTimer asynchronously.
func (t *countdownTimer) Start() {
	t.cancelled = false
	go func() {
		if t.restBeforeStart {
			t.CountdownWithSound(t.rest, "Rest", t.sound, nil)
		}
		for i := 1; i <= t.intervals-1 && !t.cancelled; i++ {
			t.CountdownWithSound(t.intervalLength, fmt.Sprintf("Interval %d/%d", i, t.intervals), t.sound, nil)
			t.CountdownWithSound(t.rest, "Rest", t.sound, nil)
		}
		t.CountdownWithSound(t.intervalLength, fmt.Sprintf("Interval %d/%d", t.intervals, t.intervals), t.sound, nil)
		t.PlaySound(t.sound, nil)
		t.doneC <- true
	}()
}

// Pause pauses the current interval countdown timer if not already paused.
func (t *countdownTimer) Pause() {
	if !t.paused {
		t.pauseC <- true
	}
}

// Resume resumes the current interval countdown timer if paused.
func (t *countdownTimer) Resume() {
	if t.paused {
		t.resumeC <- true
	}
}

// Skip skips the current interval countdownTimer. done is an optional
// channel. If done != nil it will receive a value when Cancel finishes.
func (t *countdownTimer) Skip() {
	t.skipC <- true
}

// Restart asynchronously restarts the current interval countdownTimer. done is an
// optional channel. If done != nil it will receive a value when Cancel
// finishes.
func (t *countdownTimer) Restart() {
	t.restartC <- true
}

// Cancel asynchronously cancels the countdownTimer. done is an optional channel. If
// done != nil it will receive a value when Cancel finishes.
func (t *countdownTimer) Cancel() {
	t.cancelled = true
	t.cancelC <- true
}

// Done returns a channel indicating the countdownTimer has finished.
func (t *countdownTimer) Done() <-chan bool {
	return t.doneC
}

func (t *countdownTimer) PlaySound(stream Stream, done chan<- bool) {
	if !t.cancelled {
		speaker.Clear()
		stream.Streamer.Seek(stream.StartPosition)
		speaker.Play(beep.Seq(stream.Streamer, beep.Callback(func() {
			if done != nil {
				done <- true
			}
		})))
	}
}

func (t *countdownTimer) CountdownWithSound(i Interval, name string, stream Stream, done chan<- bool) {
	if !t.cancelled {
		t.intervalName <- name
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		remaining := Interval{
			Minutes: i.Minutes,
			Seconds: i.Seconds,
		}

		for !t.cancelled {

			// Check if pause signal has been received. If it has block until
			// resume or cancel signal received.
			select {
			case <-t.pauseC:
				t.paused = true
				select {
				case <-t.resumeC:
					t.paused = false
				case <-t.cancelC:
					return
				}
			default:
			}

			select {
			case <-t.skipC:
				return
			case <-t.restartC:
				remaining.Minutes = i.Minutes
				remaining.Seconds = i.Seconds
			case <-ticker.C:
				t.timeRemaining <- remaining.String()
				switch {
				case remaining.Minutes > 0 && remaining.Seconds == 0:
					remaining.Minutes--
					remaining.Seconds = 59
				case remaining.Seconds > 0:
					remaining.Seconds--
				case remaining.Minutes == 0 && remaining.Seconds == 0:
					t.PlaySound(stream, done)
					return
				default:
				}
			}
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
