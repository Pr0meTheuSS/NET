package game

import (
	"main/geometry"
	"main/pubsub"
	"main/snake"
	"main/websnake"

	"fyne.io/fyne/v2"
)

func CreateGame(app fyne.App, username, gamename string, width, height, foodStatic, delay int32, role websnake.NodeRole) *Game {
	w := app.NewWindow("Snake")
	w.Resize(fyne.NewSize(800, 800))
	w.SetFixedSize(true)
	w.CenterOnScreen()

	thisGame := NewGame(gamename, w, []Player{}, Size{600, 600}, Size{Width: width, Height: height}, delay, foodStatic, []geometry.Position{})
	thisGame.AddMainPlayer(username, "", 0, role, websnake.PlayerType_HUMAN)

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		HandleUserInput(k, thisGame.GetMainPlayer().Snake)

		if role != websnake.NodeRole_MASTER {
			newDir := DirToNetDir[thisGame.GetMainPlayer().Snake.Dir]
			pubsub.GetGlobalPubSubService().Publish("steersend", pubsub.Message{
				Msg: &websnake.GameMessage{
					MsgSeq:     new(int64),
					SenderId:   thisGame.MainPlayerID,
					ReceiverId: new(int32),
					Type: &websnake.GameMessage_Steer{
						Steer: &websnake.GameMessage_SteerMsg{
							Direction: &newDir,
						},
					},
				},
			})
		}
	})
	w.SetOnClosed(func() { thisGame.Close() })

	w.Show()
	return thisGame
}

var DirToNetDir = map[snake.Direction]websnake.Direction{
	snake.UP:    websnake.Direction_UP,
	snake.DOWN:  websnake.Direction_DOWN,
	snake.LEFT:  websnake.Direction_LEFT,
	snake.RIGHT: websnake.Direction_RIGHT,
}
