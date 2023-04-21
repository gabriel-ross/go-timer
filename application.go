package main

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/faiface/beep"
	"github.com/gabriel-ross/timer-go/internal"
)

// Map for quick conversion of digit strings to int.
var DIGIT_MAP = genDigitStringToIntMap(100)

type Config struct {
	SpeakerSampleRate           int // Speaker ratio in HZ
	AudioBufferRatio            int // The ratio of (audio buffer size) / (sample rate)
	MaxIntervals                int
	MaxTimerMins                int
	MaxTimerSecs                int
	InitialIntervalEndSoundName string
	InitialTimerEndSoundName    string
}

type application struct {
	cnf                 Config
	guiDriver           fyne.App
	gui                 *gui
	timerConfig         *internal.Config
	timer               *internal.RepeatTimer
	speakerSampleRate   beep.SampleRate
	audioPlayer         player
	intervalFinishSound audioStream
	timerFinishSound    audioStream
	sounds              map[string]audioStream
}

func New(cnf Config, options ...func(*application)) *application {
	var err error
	newApplication := &application{
		cnf:               cnf,
		guiDriver:         app.New(),
		timerConfig:       &internal.Config{},
		speakerSampleRate: beep.SampleRate(cnf.SpeakerSampleRate),
		sounds:            map[string]audioStream{},
	}

	newApplication.audioPlayer, err = NewPlayer(newApplication.speakerSampleRate, cnf.AudioBufferRatio)
	if err != nil {
		log.Fatalf("error initializing speaker: %v", err)
	}
	newApplication.sounds["None"] = newApplication.audioPlayer.NilAudioStream(newApplication.speakerSampleRate)
	var exists bool
	newApplication.intervalFinishSound, exists = newApplication.sounds[cnf.InitialIntervalEndSoundName]
	if !exists {
		newApplication.intervalFinishSound = newApplication.sounds["None"]
		newApplication.cnf.InitialIntervalEndSoundName = "None"
	}
	newApplication.timerFinishSound, exists = newApplication.sounds[cnf.InitialTimerEndSoundName]
	if !exists {
		newApplication.timerFinishSound = newApplication.sounds["None"]
	}

	for _, option := range options {
		option(newApplication)
	}

	newApplication.gui = NewGui(newApplication)

	return newApplication
}

func WithInitialTimerConfig(cnf internal.Config) func(*application) {
	return func(a *application) {
		a.timerConfig = &cnf
	}
}

// WithAudioFiles is a functional option for configuring the sound options
// of an application. audios is a map where the key is the display name of the
// sound in the application and the value is the file path where it can be
// found. The application will crash if there are any errors registering
// the audio files.
func WithAudioFiles(audios map[string]string) func(*application) {
	return func(a *application) {
		for name, path := range audios {
			if err := a.RegisterSound(name, path); err != nil {
				log.Printf("error decoding audio file: %v\n", err)
			}
		}
	}
}

// Run runs the application.
func (a *application) Run() {
	a.gui.simpleViewWindow().ShowAndRun()
}

// OnClose handles cleanup and releases resources when an application is closed.
func (a *application) OnClose() {
	var err error
	for name, audio := range a.sounds {
		if err = audio.stream.Close(); err != nil {
			log.Printf("error closing audio %s: %v", name, err)
		}
	}
	// TODO: save current timer config
}

// RegisterSound decodes the mp3 file at path and registers it to the
// application's sound map.
func (a *application) RegisterSound(name, path string) error {
	stream, err := a.audioPlayer.NewAudioStream(path)
	if err != nil {
		return err
	}
	a.sounds[name] = stream
	return nil
}

func (a *application) runTimer() {
	a.timer = internal.NewRepeatCountdownTimer(*a.timerConfig)
	done := false

	// poll for interval name update
	go func() {
		for !done {
			select {
			case name := <-a.timer.IntervalName():
				a.gui.updateTimerName(name)
			default:
			}
		}
	}()

	// poll for time remaining updates
	go func() {
		for !done {
			select {
			case update := <-a.timer.TimeRemaining():
				a.gui.updateTimerDisplay(update)
			default:
			}
		}
	}()

	// poll for interval finished updates and play sound
	go func() {
		for !done {
			select {
			case <-a.timer.IntervalFinished():
				a.audioPlayer.PlaySound(a.speakerSampleRate, a.intervalFinishSound, nil)
			default:
			}
		}
	}()

	go func() {
		a.timer.Start()
		a.audioPlayer.PlaySound(a.speakerSampleRate, a.timerFinishSound, nil)
		done = true
		a.gui.reset()
	}()
}

func (a *application) handleTimerCancel() {
	a.timer.Cancel()
}

func (a *application) handleTimerPause() {
	a.timer.Pause()
}

func (a *application) handleTimerResume() {
	a.timer.Resume()
}

func (a *application) handleTimerSkip() {
	a.timer.Skip()
}

func (a *application) soundOptions() []string {
	opts := []string{}
	for optionKey := range a.sounds {
		opts = append(opts, optionKey)
	}
	return opts
}

func genIncrementingDigitStringSlice(start, size int) []string {
	s := []string{}
	for i := start; len(s) <= size; i++ {
		s = append(s, strconv.Itoa(i))
	}
	return s
}

func genDigitStringToIntMap(max int) map[string]int {
	m := map[string]int{}
	for i := 0; i <= max; i++ {
		m[strconv.Itoa(i)] = i
	}
	return m
}
