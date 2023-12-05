package game

import (
	"log"
	"main/pubsub"
	"main/websnake"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var gameslist = map[string]*websnake.GameAnnouncement{}
var mtx sync.RWMutex

func ChooseGame(app fyne.App, ch chan Game, username string) {
	chooseGame := app.NewWindow("Choose Game")
	chooseGame.Resize(fyne.NewSize(800, 800))
	chooseGame.SetFixedSize(true)
	chooseGame.CenterOnScreen()

	// Subscribe to event messages abount announcement
	announceSub := pubsub.Subscriber{
		EventChannel: make(chan string),
		EventHandler: func(msg pubsub.Message) {
			// log.Println("Subscriber catch message: ", msg)
			announce := msg.Msg.GetAnnouncement().GetGames()[0]
			UpdateGamesList(announce)
		},
	}

	pubsub.GetGlobalPubSubService().Subscribe("announce", announceSub)

	isChoosing := true
	chooseGame.SetOnClosed(func() { isChoosing = false })
	chooseGame.Show()
	for isChoosing {
		UpdateChooseGameWindow(app, chooseGame, ch, username)
		chooseGame.Canvas().Refresh(chooseGame.Canvas().Content())
		time.Sleep(time.Second)
	}
}

func UpdateGamesList(game *websnake.GameAnnouncement) {
	mtx.Lock()
	if _, ok := gameslist[*game.GameName]; !ok {
		gameslist[game.GetGameName()] = game
	}
	mtx.Unlock()
}

func UpdateChooseGameWindow(application fyne.App, win fyne.Window, ch chan Game, username string) {
	win.SetContent(container.NewVBox(getButtonsToConenctWithGame(application, ch, username)...))
}

func getButtonsToConenctWithGame(application fyne.App, ch chan Game, username string) []fyne.CanvasObject {
	btns := []fyne.CanvasObject{}
	mtx.RLock()
	for k, v := range gameslist {
		btns = append(btns, widget.NewButton(k, func() {
			log.Println("Choose to connection")

			game := CreateGame(application, username, *v.GameName, *v.Config.Width, *v.Config.Height, *v.Config.FoodStatic, *v.Config.StateDelayMs, websnake.NodeRole_NORMAL)
			go game.ConnectToTheGame(v)
			ch <- *game
		}))
	}
	mtx.RUnlock()

	return btns
}
