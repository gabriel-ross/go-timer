package timer

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gabriel-ross/timer-go/internal"
)

var DIGIT_MAP = genDigitStringToIntMap(100)

type Config struct {
	maxIntervals int
}

type application struct {
	gui                  fyne.App
	cnf                  Config
	timerConfig          *internal.Config
	timerNameDisplay     binding.String
	timeRemainingDisplay binding.String
}

func (a *application) configureSimpleViewWindow() fyne.Window {
	w := a.gui.NewWindow("simple timer")
	timerName := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	timerName.Bind(a.timerNameDisplay)
	timeRemaining := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	timeRemaining.Bind(a.timeRemainingDisplay)
	displayVBox := container.New(layout.NewVBoxLayout(), timerName, timeRemaining)

	intervalsSelect := widget.NewSelect(genIncrementingDigitStringSlice(1, a.cnf.maxIntervals), func(s string) {

	})
	resumeButton := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {})
	pauseButton := widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() {})
	skipButton := widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), func() {})
	stopButton := widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() {})
	buttonGrid := container.New(layout.NewGridLayout(2), resumeButton, pauseButton, stopButton, skipButton)

	windowVBox := container.New(layout.NewVBoxLayout(), displayVBox, layout.NewSpacer(), layout.NewSpacer(), buttonGrid)
	w.SetContent(windowVBox)

	return w
}

func (a *application) configureMenuWindow() fyne.Window {
	w := a.gui.NewWindow("menu")

	numIntervalOptions := widget.NewSelect(a.numIntervalOptions, func(s string) {
		a.timerCnf.Intervals = a.numIntervalDigitMap[s]
	})
	numIntervalOptions.SetSelectedIndex(0)

	soundOptions := widget.NewSelect(genMapListOptions(SOUND_PATHS), func(s string) {})
	soundOptions.SetSelectedIndex(0)

	intervalDurationMins := &widget.Select{
		Options:     a.numIntervalOptions,
		PlaceHolder: "MM",
		OnChanged: func(s string) {
			a.timerCnf.IntervalLength.Minutes = a.numIntervalDigitMap[s]
		},
	}
	intervalDurationSecs := &widget.Select{
		Options:     a.numIntervalOptions,
		PlaceHolder: "SS",
		OnChanged: func(s string) {
			a.timerCnf.IntervalLength.Seconds = a.numIntervalDigitMap[s]
		},
	}
	intervalLength := container.New(layout.NewHBoxLayout(), intervalDurationMins, widget.NewLabel(":"), intervalDurationSecs)

	restDurationMins := &widget.Select{
		Options:     a.numIntervalOptions,
		PlaceHolder: "MM",
		OnChanged: func(s string) {
			a.timerCnf.Rest.Minutes = a.numIntervalDigitMap[s]
		},
	}
	restDurationSecs := &widget.Select{
		Options:     a.numIntervalOptions,
		PlaceHolder: "SS",
		OnChanged: func(s string) {
			a.timerCnf.Rest.Seconds = a.numIntervalDigitMap[s]
		},
	}
	restLength := container.New(layout.NewHBoxLayout(), restDurationMins, widget.NewLabel(":"), restDurationSecs)

	restBeforeStart := widget.NewCheck("", func(checked bool) {
		a.timerCnf.RestBeforeStart = checked
	})

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Rest", Widget: restLength},
			{Text: "Interval", Widget: intervalLength},
			{Text: "# Intervals", Widget: numIntervalOptions},
			{Text: "Sound", Widget: soundOptions},
			{Text: "Rest before start", Widget: restBeforeStart},
		},
	}
	start := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {})

	w.SetContent(container.NewBorder(nil, start, nil, nil, form))
	return w
}

func genIncrementingDigitStringSlice(start, size int) []string {
	s := []string{}
	for i := start; len(s) <= size; i++ {
		s = append(s, strconv.Itoa(i))
	}
	return s
}

func genMapListOptions(options map[string]string) []string {
	opts := []string{}
	for name := range options {
		opts = append(opts, name)
	}
	return opts
}

func genDigitStringToIntMap(max int) map[string]int {
	m := map[string]int{}
	for i := 0; i <= max; i++ {
		str, _ := strconv.Atoi(i)
		m[str] = i
	}
	return m
}

func genDigitMapFromSlice(s []string) map[string]int {
	dMap := map[string]int{}
	startingVal, _ := strconv.Atoi(s[0])
	for idx, val := range s {
		dMap[val] = idx + startingVal
	}
	return dMap
}

func validateDigitString(s string) error {
	_, err := strconv.Atoi(s)
	return err
}
