package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {

	app := app.New()

	w := app.NewWindow("11111")
	w.Resize(fyne.NewSize(500, 700))

	can := widget.NewButton("asd", func() {

		dialog.ShowCustom("wre", "wer", widget.NewLabel("sfsdfsdfsdf"), w)
	})
	can.Resize(fyne.NewSize(100, 100))

	w.SetContent(container.NewWithoutLayout(can))
	w.ShowAndRun()
}
