package main

import (
	"github.com/gabriel-ross/timer-go"
)

var (
	MAX_INTERVALS = 99
	SOUND_PATHS   = map[string]string{
		"ding": "sounds/iphone-ding-sound.mp3",
	}
	WINDOW_WIDTH  float32 = 300
	WINDOW_HEIGHT float32 = 800

	default_timer_config = timer.Config{
		Intervals: 1,
		IntervalLength: timer.Interval{
			Minutes: 0,
			Seconds: 30,
		},
		Rest: timer.Interval{
			Minutes: 0,
			Seconds: 5,
		},
		SoundPath:       SOUND_PATHS["ding"],
		RestBeforeStart: false,
	}
)

func main() {
	app := timer.New()
	app.Start()
}
