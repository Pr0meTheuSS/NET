package frames

import (
	"main/game"
	"main/webnodes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func InitHelloWinContent(application fyne.App) *fyne.Container {
	usernameEntry := widget.NewEntry()

	connectToTheGameButton := widget.NewButton("Подключиться к существующей игре", func() {
		webNode := webnodes.NewWebSnakeNormalNode(usernameEntry.Text)
		webNode.Run()
		game.ChooseGame(application, usernameEntry.Text)
	})
	createTheGameButton := widget.NewButton("Создать новую игру", func() {
		createAndShowSetGameConfigFrame(application, usernameEntry.Text)
	})

	hello := widget.NewLabel("Hello, what is your name?")

	return container.NewVBox(
		hello,
		usernameEntry,
		connectToTheGameButton,
		createTheGameButton,
	)
}
