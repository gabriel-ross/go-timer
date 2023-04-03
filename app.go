package timer

import (
	"strconv"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	MAX_INTERVALS = 99
	SOUND_PATHS   = map[string]string{
		"ding": "sounds/iphone-ding-sound.mp3",
	}
	WINDOW_WIDTH  float32 = 300
	WINDOW_HEIGHT float32 = 800

	default_timer_config = Config{
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

type app struct {
	gui                 fyne.App
	menuWindow          fyne.Window
	timerWindow         fyne.Window
	intervalNameDisplay binding.String
	timerDisplay        binding.String
	timer               *countdownTimer
	timerCnf            Config
}

func New() *app {
	a := &app{
		gui:                 fyneApp.New(),
		timer:               &countdownTimer{},
		timerCnf:            default_timer_config,
		intervalNameDisplay: binding.NewString(),
		timerDisplay:        binding.NewString(),
	}
	a.intervalNameDisplay.Set("Starting")
	a.timerDisplay.Set("00:00")
	a.menuWindow = a.configureMenuWindow()
	a.timerWindow = a.configureTimerWindow()

	return a
}

func (a *app) Start() {
	a.menuWindow.Show()
	a.gui.Run()
}

func (a *app) configureMenuWindow() fyne.Window {
	w := a.gui.NewWindow("menu")

	numIntervalOptions := widget.NewSelect(genNumIntervalOptions(), func(s string) {})
	soundOptions := widget.NewSelect(genSoundOptions(), func(s string) {})

	restMinsText := binding.NewString()
	restSecsText := binding.NewString()
	restMins := binding.StringToInt(restMinsText)
	restSecs := binding.StringToInt(restSecsText)
	restMinsEntry := widget.NewEntryWithData(restMinsText)
	restSecsEntry := widget.NewEntryWithData(restSecsText)
	restMinsEntry.Wrapping = fyne.TextWrapOff
	restSecsEntry.Wrapping = fyne.TextWrapOff
	restMinsEntry.Validator = validateDigitString
	restSecsEntry.Validator = validateDigitString
	restMinsEntry.PlaceHolder = "MM"
	restSecsEntry.PlaceHolder = "SS"
	restLength := container.New(layout.NewHBoxLayout(), restMinsEntry, widget.NewLabel(":"), restSecsEntry)

	intervalMinsText := binding.NewString()
	intervalSecsText := binding.NewString()
	intervalMins := binding.StringToInt(intervalMinsText)
	intervalSecs := binding.StringToInt(intervalSecsText)
	intervalMinsEntry := widget.NewEntryWithData(intervalMinsText)
	intervalSecsEntry := widget.NewEntryWithData(intervalSecsText)
	intervalMinsEntry.Wrapping = fyne.TextWrapOff
	intervalSecsEntry.Wrapping = fyne.TextWrapOff
	intervalMinsEntry.Validator = validateDigitString
	intervalSecsEntry.Validator = validateDigitString
	intervalMinsEntry.PlaceHolder = "MM"
	intervalSecsEntry.PlaceHolder = "SS"
	intervalLength := container.New(layout.NewHBoxLayout(), intervalMinsEntry, widget.NewLabel(":"), intervalSecsEntry)

	restBeforeStart := widget.NewCheck("", nil)

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
		a.timerCnf.Rest.Minutes, _ = restMins.Get()
		a.timerCnf.Rest.Seconds, _ = restSecs.Get()
		a.timerCnf.IntervalLength.Minutes, _ = intervalMins.Get()
		a.timerCnf.IntervalLength.Seconds, _ = intervalSecs.Get()
		a.timerCnf.Intervals, _ = strconv.Atoi(numIntervalOptions.Selected)
		a.timerCnf.SoundPath = SOUND_PATHS[soundOptions.Selected]
		a.timerCnf.RestBeforeStart = restBeforeStart.Checked
		a.runAndDisplayTimer()
	})

	w.SetContent(container.NewBorder(nil, start, nil, nil, form))
	return w
}

func (a *app) runAndDisplayTimer() {
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

func (a *app) configureTimerWindow() fyne.Window {
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

func (a *app) swapWindow(cur fyne.Window, next fyne.Window) {
	next.Show()
	cur.Close()
}

func (a *app) handlePauseButton() {
	a.timer.Pause()
}

func (a *app) handleResumeButton() {
	a.timer.Resume()
}

func (a *app) handleSkipButton() {
	a.timer.Skip()
}

func (a *app) handleRestartButton() {
	a.timer.Restart()
}

func (a *app) handleTimerStop() {
	a.timer.Cancel()
	a.menuWindow.Show()
}

func genNumIntervalOptions() []string {
	opts := []string{}
	for i := 1; i <= 99; i++ {
		opts = append(opts, strconv.Itoa(i))
	}
	return opts
}

func genSoundOptions() []string {
	opts := []string{}
	for name := range SOUND_PATHS {
		opts = append(opts, name)
	}
	return opts
}

func validateDigitString(s string) error {
	_, err := strconv.Atoi(s)
	return err
}
