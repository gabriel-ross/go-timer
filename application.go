package timer

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/gabriel-ross/timer-go/internal"
)

var DIGIT_MAP = genDigitStringToIntMap(100)

type Config struct {
	MaxIntervals int
	MaxTimerMins int
	MaxTimerSecs int
}

type application struct {
	cnf                 Config
	guiDriver           fyne.App
	gui                 *gui
	timerConfig         *internal.Config
	timer               *internal.RepeatTimer
	audioPlayer         player
	intervalFinishSound audioStream
	timerFinishSound    audioStream
	sounds              map[string]audioStream
}

func New(cnf Config, options ...func(*application)) *application {
	newApplication := &application{
		cnf:         cnf,
		guiDriver:   app.New(),
		timerConfig: &internal.Config{},
	}
	newApplication.gui = NewGui(newApplication)

	return newApplication
}

// WithAudioFiles is a functional option for configuring the sound options
// of an application. audios is a map where the key is the display name of the
// sound in the application and the value is the file path where it can be
// found. The application will crash if there are any errors registering
// the audio files.
func WithAudioFiles(audios map[string]string) func(*application) {
	return func(a *application) {
		for name, path := range audios {
			a.MustRegisterSound(name, path)
		}
	}
}

func (a *application) Start() {
	w := a.gui.simpleViewWindow()
	w.Show()
	a.guiDriver.Run()
}

func (a *application) MustRegisterSound(name, path string) {
	stream, err := a.audioPlayer.newAudioStream(path)
	if err != nil {
		log.Fatalf("error registering audio: %v", err)
	}
	a.sounds[name] = stream
}

func (a *application) startTimer() {
	a.timer = internal.NewRepeatCountdownTimer(*a.timerConfig)
	done := make(chan bool)

	go func() {
		go a.pollForIntervalNameUpdates()
		go a.pollForTimeRemainingUpdates()
		go a.pollForIntervalFinishedAndPlaySound(a.intervalFinishSound)
		<-done
	}()

	go func() {
		a.timer.Start()
		a.audioPlayer.playSound(a.timerFinishSound, nil)
		done <- true
		a.gui.reset()
	}()
}

func (a *application) pollForIntervalNameUpdates() {
	for {
		a.gui.updateTimerName(<-a.timer.IntervalName())
	}
}

func (a *application) pollForTimeRemainingUpdates() {
	for {
		a.gui.updateTimerDisplay(<-a.timer.TimeRemaining())
	}
}

func (a *application) pollForIntervalFinishedAndPlaySound(audio audioStream) {
	for {
		<-a.timer.IntervalFinished()
		a.audioPlayer.playSound(audio, nil)
	}
}

func (a *application) handleTimerCancel() {
	a.timer.Cancel()
	a.gui.reset()
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
	for optionKey, _ := range a.sounds {
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
