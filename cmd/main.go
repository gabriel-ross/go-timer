package main

import (
	"time"

	"github.com/gabriel-ross/timer-go"
)

var TIMEOUT = 2 * time.Second

func main() {
	// f, err := os.Open("sounds/iphone-ding-sound.mp3")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// streamer, format, err := mp3.Decode(f)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer streamer.Close()

	// speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	// done := make(chan bool)
	// speaker.Play(beep.Seq(streamer, beep.Callback(func() {
	// 	done <- true
	// })))
	// select {
	// case <-done:
	// 	print("done")
	// case <-time.After(TIMEOUT):
	// 	print("timed out")
	// }

	// ticker := time.NewTicker(time.Second)
	// defer ticker.Stop()
	// done := make(chan bool)
	// go func() {
	// 	time.Sleep(10 * time.Second)
	// 	done <- true
	// }()
	// for {
	// 	select {
	// 	case <-done:
	// 		fmt.Println("Done!")
	// 		return
	// 	case t := <-ticker.C:
	// 		fmt.Println("Current time: ", t)
	// 	}
	// }

	t := timer.Timer{
		Intervals: 5,
		IntervalLength: timer.Interval{
			Minutes: 0,
			Seconds: 5,
		},
		Rest: timer.Interval{
			Minutes: 0,
			Seconds: 5,
		},
		IntervalSound:   "",
		RestSound:       "",
		RestBeforeStart: false,
	}

	t.Start()
}
