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
}
