package main

import (
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")
	myWindow.SetContent(widget.NewLabel("Hello\nWait 5 seconds..."))
	myWindow.Resize(fyne.NewSize(800, 600))

	go showAnother(myApp, myWindow)
	myWindow.ShowAndRun()
}

func showAnother(a fyne.App, w fyne.Window) {
	time.Sleep(time.Second * 5)

	win := a.NewWindow("Shown later")
	win.SetContent(widget.NewLabel("Waiting for 3 seconds to hide"))
	win.Resize(fyne.NewSize(200, 200))
	win.Show()

	time.Sleep(time.Second * 3)
	win.Hide()

  w.SetContent(widget.NewLabel(
    "Waiting 2 seconds to close the second window...",
  ))
	time.Sleep(time.Second * 2)
  w.SetContent(widget.NewLabel(
    "Waiting 2 seconds to close the second window...\nSecond window closed",
  ))
  win.Close()
}
