package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
  "fyne.io/fyne/v2"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Widget")

  myWindow.Resize(fyne.NewSize(150, 50))
	myWindow.SetContent(widget.NewEntry())
	myWindow.ShowAndRun()
}
