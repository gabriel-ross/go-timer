package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("demo")

	timerName := widget.NewLabelWithStyle("Interval 1/2", fyne.TextAlignCenter, fyne.TextStyle{})
	timeRemaining := widget.NewLabelWithStyle("00:00", fyne.TextAlignCenter, fyne.TextStyle{})
	displayVBox := container.New(layout.NewVBoxLayout(), timerName, timeRemaining)

	resumeButton := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {})
	pauseButton := widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() {})
	skipButton := widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), func() {})
	stopButton := widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() {})
	buttonGrid := container.New(layout.NewGridLayout(2), resumeButton, pauseButton, stopButton, skipButton)

	windowVBox := container.New(layout.NewVBoxLayout(), displayVBox, layout.NewSpacer(), layout.NewSpacer(), buttonGrid)
	w.SetContent(windowVBox)

	w.ShowAndRun()
}
