package frames

import (
	"log"
	"main/game"
	"main/webnodes"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func InitHelloWinContent(application fyne.App) *fyne.Container {
	usernameEntry := widget.NewEntry()

	ch := make(chan game.Game)

	connectToTheGameButton := widget.NewButton("Подключиться к существующей игре", func() {
		webNode := webnodes.NewWebSnakeNormalNode()
		webNode.Run()

		go game.ChooseGame(application, ch)

		g := <-ch
		log.Println("After generation the game")
		webNode.SetGame(&g)
	})

	createTheGameButton := widget.NewButton("Создать новую игру", func() {
		createAndShowSetGameConfigFrame(application, usernameEntry.Text, ch)

		g := <-ch
		log.Println("After channel reading", g)

		node := webnodes.NewWebNode(&g)
		go g.MainLoop()
		node.RunLikeMaster()
	})

	hello := widget.NewLabel("Hello, what is your name?")

	return container.NewVBox(
		hello,
		usernameEntry,
		connectToTheGameButton,
		createTheGameButton,
	)
}
