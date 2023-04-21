package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	DEFAULT_TIMER_NAME    = "Current interval name"
	DEFAULT_TIMER_DISPLAY = "00:00"
)

type gui struct {
	application         *application
	intervals           *widget.Select
	sounds              *widget.Select
	intervalDurationMin *widget.Select
	intervalDurationSec *widget.Select
	restDurationMin     *widget.Select
	restDurationSec     *widget.Select
	restBeforeStart     *widget.Check
	timerName           *canvas.Text
	timeRemaining       *canvas.Text
	stopButton          *widget.Button
	pauseButton         *widget.Button
	startResumeButton   *widget.Button
	skipButton          *widget.Button
}

func NewGui(app *application) *gui {
	newGui := &gui{
		application: app,
	}

	newGui.timerName = canvas.NewText(DEFAULT_TIMER_NAME, color.Black)
	newGui.timerName.Alignment = fyne.TextAlignCenter
	newGui.timeRemaining = canvas.NewText(DEFAULT_TIMER_DISPLAY, color.Gray16{3})
	newGui.timeRemaining.TextSize = 50
	newGui.timeRemaining.Alignment = fyne.TextAlignCenter

	return newGui
}

func (g *gui) simpleViewWindow() fyne.Window {
	w := g.application.guiDriver.NewWindow("simple view")

	displayVBox := container.New(layout.NewVBoxLayout(), g.timerName, g.timeRemaining)

	intervalsLabel := g.newCenteredText("# of Intervals", color.Black)
	soundsLabel := g.newCenteredText("Sound", color.Black)
	intervalLabel := g.newCenteredText("Interval", color.Black)
	restLabel := g.newCenteredText("Rest", color.Black)
	restBeforeStartLabel := g.newCenteredText("Rest before start", color.Black)
	g.intervals = widget.NewSelect(genIncrementingDigitStringSlice(1, g.application.cnf.MaxIntervals), g.handleIntervalsSelect)
	g.sounds = widget.NewSelect(g.application.soundOptions(), g.handleSoundSelect)
	g.sounds.SetSelected(g.application.cnf.InitialIntervalEndSoundName)
	g.intervalDurationMin = &widget.Select{
		Options:     genIncrementingDigitStringSlice(0, g.application.cnf.MaxTimerMins),
		PlaceHolder: "MM",
		OnChanged:   g.handleIntervalMinuteSelect,
	}
	g.intervalDurationSec = &widget.Select{
		Options:     genIncrementingDigitStringSlice(0, g.application.cnf.MaxTimerSecs),
		PlaceHolder: "SS",
		OnChanged:   g.handleIntervalSecondSelect,
	}
	g.restDurationMin = &widget.Select{
		Options:     genIncrementingDigitStringSlice(0, g.application.cnf.MaxTimerMins),
		PlaceHolder: "MM",
		OnChanged:   g.handleRestMinuteSelect,
	}
	g.restDurationSec = &widget.Select{
		Options:     genIncrementingDigitStringSlice(0, g.application.cnf.MaxTimerSecs),
		PlaceHolder: "SS",
		OnChanged:   g.handleRestSecondSelect,
	}
	interval := container.New(layout.NewHBoxLayout(), g.intervalDurationMin, widget.NewLabel(":"), g.intervalDurationSec)
	rest := container.New(layout.NewHBoxLayout(), g.restDurationMin, widget.NewLabel(":"), g.restDurationSec)
	g.restBeforeStart = widget.NewCheck("", g.handleRestBeforeStartChecked)

	settings := container.New(layout.NewGridLayout(2),
		intervalsLabel, g.intervals,
		intervalLabel, interval,
		restLabel, rest,
		restBeforeStartLabel, g.restBeforeStart,
		soundsLabel, g.sounds)

	g.stopButton = widget.NewButtonWithIcon("", theme.MediaStopIcon(), g.handleStopButtonTap)
	g.pauseButton = widget.NewButtonWithIcon("", theme.MediaPauseIcon(), g.handlePauseButtonTap)
	g.startResumeButton = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), g.handleStartButtonTap)
	g.skipButton = widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), g.handleSkipButtonTap)
	g.startResumeButton.Enable()
	g.pauseButton.Disable()
	g.stopButton.Disable()
	g.skipButton.Disable()
	buttonGrid := container.New(layout.NewGridLayout(2), g.skipButton, g.pauseButton, g.stopButton, g.startResumeButton)

	windowVBox := container.New(layout.NewVBoxLayout(), displayVBox, layout.NewSpacer(), settings, layout.NewSpacer(), buttonGrid)
	w.SetContent(windowVBox)

	return w
}

func (g *gui) reset() {

	g.updateTimerName(DEFAULT_TIMER_NAME)
	g.updateTimerDisplay(DEFAULT_TIMER_DISPLAY)
	g.intervals.Enable()
	g.intervalDurationMin.Enable()
	g.intervalDurationSec.Enable()
	g.restDurationMin.Enable()
	g.restDurationSec.Enable()
	g.restBeforeStart.Enable()
	g.startResumeButton.Enable()
	g.pauseButton.Disable()
	g.stopButton.Disable()
	g.skipButton.Disable()
	g.startResumeButton.OnTapped = g.handleStartButtonTap
}

func (g *gui) updateTimerName(text string) {
	g.timerName.Text = text
	g.timerName.Refresh()
}

func (g *gui) updateTimerDisplay(text string) {
	g.timeRemaining.Text = text
	g.timeRemaining.Refresh()
}

func (g *gui) handleIntervalsSelect(s string) {
	g.application.timerConfig.Intervals = DIGIT_MAP[s]
}

func (g *gui) handleSoundSelect(s string) {
	g.application.intervalFinishSound = g.application.sounds[s]
	g.application.timerFinishSound = g.application.sounds[s]
}

func (g *gui) handleIntervalMinuteSelect(s string) {
	g.application.timerConfig.IntervalMinutes = int64(DIGIT_MAP[s])
}

func (g *gui) handleIntervalSecondSelect(s string) {
	g.application.timerConfig.IntervalSeconds = int64(DIGIT_MAP[s])
}

func (g *gui) handleRestMinuteSelect(s string) {
	g.application.timerConfig.RestMinutes = int64(DIGIT_MAP[s])
}

func (g *gui) handleRestSecondSelect(s string) {
	g.application.timerConfig.RestSeconds = int64(DIGIT_MAP[s])
}

func (g *gui) handleRestBeforeStartChecked(checked bool) {
	g.application.timerConfig.RestBeforeStart = checked
}

func (g *gui) handleStartButtonTap() {
	g.intervals.Disable()
	g.intervalDurationMin.Disable()
	g.intervalDurationSec.Disable()
	g.restDurationMin.Disable()
	g.restDurationSec.Disable()
	g.restBeforeStart.Disable()

	g.pauseButton.Enable()
	g.stopButton.Enable()
	g.startResumeButton.Disable()
	g.skipButton.Enable()
	g.startResumeButton.OnTapped = g.handleResumeButtonTap
	g.application.runTimer()
}

func (g *gui) handlePauseButtonTap() {
	g.pauseButton.Disable()
	g.startResumeButton.Enable()
	g.application.handleTimerPause()
}

func (g *gui) handleResumeButtonTap() {
	g.pauseButton.Enable()
	g.application.handleTimerResume()
}

func (g *gui) handleSkipButtonTap() {
	g.application.handleTimerSkip()
}

func (g *gui) handleStopButtonTap() {
	g.application.handleTimerCancel()
	g.reset()
}

func (g gui) newCenteredText(text string, color color.Color) *canvas.Text {
	newText := canvas.NewText(text, color)
	newText.Alignment = fyne.TextAlignCenter
	return newText
}
