package main

import (
	"main/frames"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	snakeApp := app.New()
	snakeMainWindow := snakeApp.NewWindow("Snake web")
	snakeMainWindow.SetMaster()

	snakeMainWindow.SetContent(frames.InitHelloWinContent(snakeApp))
	snakeMainWindow.Resize(fyne.NewSize(600.0, 600.0))

	snakeMainWindow.ShowAndRun()
}
