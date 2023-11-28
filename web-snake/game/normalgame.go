package game

import (
	"log"
	"main/pubsub"
	"main/websnake"
	"net"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var gameslist = map[string]*websnake.GameAnnouncement{}
var mtx sync.RWMutex
var uname = ""

func ChooseGame(app fyne.App, username string) {
	uname = username
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

	// Subscribe to event messages abount ack
	ackSub := pubsub.Subscriber{
		EventChannel: make(chan string),
		EventHandler: func(msg pubsub.Message) {
			log.Println("Receive ack and set main player id as receiver id from message")
			thisGame.mainPlayerId = *msg.Msg.ReceiverId
		},
	}

	pubsub.GetGlobalPubSubService().Subscribe("ack", ackSub)

	isChoosing := true
	chooseGame.SetOnClosed(func() { isChoosing = false })
	chooseGame.Show()
	for isChoosing {
		UpdateChooseGameWindow(app, chooseGame)
		chooseGame.Canvas().Refresh(chooseGame.Canvas().Content())
		time.Sleep(time.Second)
	}
}

func UpdateGamesList(game *websnake.GameAnnouncement) {
	if _, ok := gameslist[*game.GameName]; !ok {
		mtx.Lock()
		gameslist[game.GetGameName()] = game
		mtx.Unlock()
	}
}

func UpdateChooseGameWindow(application fyne.App, win fyne.Window) {
	win.SetContent(container.NewVBox(getButtonsToConenctWithGame(application)...))
}

var thisGame = &Game{}

func getButtonsToConenctWithGame(application fyne.App) []fyne.CanvasObject {
	btns := []fyne.CanvasObject{}
	mtx.RLock()
	for k, v := range gameslist {
		btns = append(btns, widget.NewButton(k, func() {
			log.Println("Choose to connection")

			masterUDPAddr := getMasterAddressFromAnnounce(v)
			pubsub.GetGlobalPubSubService().Publish("connection", pubsub.Message{
				Msg: &websnake.GameMessage{
					MsgSeq:     new(int64),
					SenderId:   new(int32),
					ReceiverId: new(int32),
					Type: &websnake.GameMessage_Announcement{
						Announcement: &websnake.GameMessage_AnnouncementMsg{
							Games: []*websnake.GameAnnouncement{v},
						},
					},
				},
				From: nil,
				To:   masterUDPAddr,
			})
			thisGame = CreateGame(application, *v.GameName, *v.GameName, *v.Config.Width, *v.Config.Height, *v.Config.FoodStatic, *v.Config.StateDelayMs)
		}))
	}
	mtx.RUnlock()

	return btns
}

func getMasterAddressFromAnnounce(announce *websnake.GameAnnouncement) *net.UDPAddr {
	players := announce.Players.Players
	for _, player := range players {
		if *player.Role == *websnake.NodeRole_MASTER.Enum() {
			port := int(*player.Port)
			addr := net.UDPAddr{
				IP:   net.ParseIP(*player.IpAddress),
				Port: port,
				Zone: "",
			}

			return &addr
		}
	}

	return nil
}
