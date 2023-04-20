package main

import "github.com/gabriel-ross/timer-go"

var (
	IPHONE_SPEAKER_SAMPLE_RATE_HZ = 48000
	DEFAULT_AUDIO_BUFFER_RATIO    = 5
	AUDIO_FILES                   = map[string]string{
		"Ding": "./assets/audio/iphone-ding-sound.mp3",
	}
)

func main() {
	a := timer.New(timer.Config{
		SpeakerSampleRate:           IPHONE_SPEAKER_SAMPLE_RATE_HZ,
		AudioBufferRatio:            DEFAULT_AUDIO_BUFFER_RATIO,
		MaxIntervals:                99,
		MaxTimerMins:                99,
		MaxTimerSecs:                59,
		InitialIntervalEndSoundName: "Ding",
		InitialTimerEndSoundName:    "Chime",
	}, timer.WithAudioFiles(AUDIO_FILES))
	a.Run()

	// t := internal.NewRepeatCountdownTimer(internal.Config{
	// 	Intervals:       10,
	// 	IntervalMinutes: 0,
	// 	IntervalSeconds: 10,
	// 	RestEnabled:     true,
	// 	RestMinutes:     0,
	// 	RestSeconds:     5,
	// 	RestBeforeStart: false,
	// })
	// done := false

	// go func() {
	// 	// poll for time remaining updates
	// 	go func() {
	// 		for !done {
	// 			select {
	// 			case timeRemaining := <-t.TimeRemaining():
	// 				fmt.Println("time remaining: ", timeRemaining)
	// 			default:
	// 			}
	// 		}
	// 	}()

	// 	// poll for interval name updates
	// 	go func() {
	// 		for !done {
	// 			select {
	// 			case intervalName := <-t.IntervalName():
	// 				fmt.Println(intervalName)
	// 			default:
	// 			}
	// 		}
	// 	}()
	// }()

	// go func() {
	// 	time.Sleep(3 * time.Second)
	// 	// fmt.Println("cancelling...")
	// 	// t.Cancel()
	// 	fmt.Println("pausing...")
	// 	t.Pause()
	// 	time.Sleep(3 * time.Second)
	// 	fmt.Println("Resuming")
	// 	t.Resume()
	// }()

	// t.Start()

	// done = true

	// time.Sleep(3 * time.Second)
}
