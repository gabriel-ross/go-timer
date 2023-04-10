package timer

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	MAX_DIGIT_OPTIONS = 99
	DIGITS            = genDigitListOptions(MAX_DIGIT_OPTIONS)
	DIGIT_MAP         = genDigitMap(DIGITS)
	SOUND_PATHS       = map[string]string{
		"ding": "sounds/iphone-ding-sound.mp3",
	}
	WINDOW_WIDTH  float32 = 300
	WINDOW_HEIGHT float32 = 800

	default_timer_config = timerConfig{
		Intervals: 1,
		IntervalLength: Interval{
			Minutes: 0,
			Seconds: 30,
		},
		Rest: Interval{
			Minutes: 0,
			Seconds: 5,
		},
		SoundPath:       SOUND_PATHS["ding"],
		RestBeforeStart: false,
	}
)

type Config struct {
	MaxIntervals        int
	MaxIntervalDuration Interval
}

type application struct {
	gui                 fyne.App
	menuWindow          fyne.Window
	timerWindow         fyne.Window
	intervalNameDisplay binding.String
	timerDisplay        binding.String
	timer               *countdownTimer
	timerCnf            timerConfig
	cnf                 Config
	numIntervalOptions  []string
	numIntervalDigitMap map[string]int
}

func New(cnf Config) *application {
	a := &application{
		gui:                 app.New(),
		timer:               &countdownTimer{},
		timerCnf:            default_timer_config,
		intervalNameDisplay: binding.NewString(),
		timerDisplay:        binding.NewString(),
	}
	a.numIntervalOptions = genDigitListOptions(cnf.MaxIntervals)
	a.numIntervalDigitMap = genDigitMap(a.numIntervalOptions)
	a.intervalNameDisplay.Set("Starting")
	a.timerDisplay.Set("00:00")
	a.menuWindow = a.configureMenuWindow()
	a.timerWindow = a.configureTimerWindow()

	return a
}

func (a *application) Start() {
	a.menuWindow.Show()
	a.gui.Run()
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
	start := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		a.runTimerAndDisplay()
	})

	w.SetContent(container.NewBorder(nil, start, nil, nil, form))
	return w
}

func (a *application) runTimerAndDisplay() {
	t, err := NewTimer(a.timerCnf)
	if err != nil {
		dialog.ShowError(err, a.menuWindow)
		return
	}
	a.timer = t
	a.timerWindow = a.configureTimerWindow()
	a.timer.Start()
	a.swapWindow(a.menuWindow, a.timerWindow)

	// Poll for timer display updates and interval name display updates
	go func() {
		for {
			a.timerDisplay.Set(<-a.timer.timeRemaining)
			select {
			case intervalName := <-a.timer.intervalName:
				a.intervalNameDisplay.Set(intervalName)
			default:
			}
		}
	}()
	<-a.timer.Done()
	a.swapWindow(a.timerWindow, a.menuWindow)
}

func (a *application) configureTimerWindow() fyne.Window {
	w := a.gui.NewWindow("timer")

	intervalName := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})
	timerDisplay := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	intervalName.Bind(a.intervalNameDisplay)
	timerDisplay.Bind(a.timerDisplay)
	timerVBox := container.New(layout.NewVBoxLayout(), intervalName, timerDisplay)

	pauseButton := widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() {})
	resumeButton := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {})
	resumeButton.Disable()

	pauseButton.OnTapped = func() {
		pauseButton.Disable()
		resumeButton.Enable()
		a.timer.Pause()
	}
	resumeButton.OnTapped = func() {
		pauseButton.Enable()
		resumeButton.Disable()
		a.timer.Resume()
	}

	skipButton := widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), a.handleSkipButton)
	restartButton := widget.NewButtonWithIcon("", theme.MediaReplayIcon(), a.handleRestartButton)
	stopButton := widget.NewButtonWithIcon("", theme.MediaStopIcon(), a.handleTimerStop)
	vBox := container.New(layout.NewVBoxLayout(), timerVBox, pauseButton, resumeButton, skipButton, restartButton, stopButton)

	w.SetContent(vBox)
	return w
}

func (a *application) swapWindow(cur fyne.Window, next fyne.Window) {
	next.Show()
	cur.Close()
}

func (a *application) handlePauseButton() {
	a.timer.Pause()
}

func (a *application) handleResumeButton() {
	a.timer.Resume()
}

func (a *application) handleSkipButton() {
	a.timer.Skip()
}

func (a *application) handleRestartButton() {
	a.timer.Restart()
}

func (a *application) handleTimerStop() {
	a.timer.Cancel()
	a.menuWindow.Show()
}

func genDigitListOptions(max int) []string {
	opts := []string{}
	for i := 1; i <= max; i++ {
		opts = append(opts, strconv.Itoa(i))
	}
	return opts
}

func genMapListOptions(options map[string]string) []string {
	opts := []string{}
	for name := range options {
		opts = append(opts, name)
	}
	return opts
}

func genDigitMap(s []string) map[string]int {
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
