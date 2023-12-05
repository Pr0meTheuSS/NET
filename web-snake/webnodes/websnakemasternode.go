package webnodes

import (
	"fmt"
	"log"
	"main/game"
	"main/pubsub"
	"main/snake"
	"main/websnake"
	"net"
	"os"
	"time"

	"google.golang.org/protobuf/proto"
)

type WebSnakeMasterNode struct {
	conn *net.UDPConn
	game *game.Game
}

func (w *WebSnakeMasterNode) Run() {
	go w.SendMultiAnnouncment()
	go w.ListenAndServe()
}

func NewWebSnakeMasterNode(game *game.Game) *WebSnakeMasterNode {
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Fatal(err)
	}

	node := &WebSnakeMasterNode{
		conn: conn,
		game: game,
	}

	ackSender := pubsub.Subscriber{
		EventChannel: make(chan string),
		EventHandler: func(msg pubsub.Message) {
			node.sendAck(msg)
		},
	}

	pubsub.GetGlobalPubSubService().Subscribe("sendack", ackSender)

	gameStateSender := pubsub.Subscriber{
		EventChannel: make(chan string),
		EventHandler: func(msg pubsub.Message) {
			node.sendGameState(msg)
		},
	}

	pubsub.GetGlobalPubSubService().Subscribe("sendgamestate", gameStateSender)

	return node
}

var DirToNetDir = map[snake.Direction]websnake.Direction{
	snake.UP:    websnake.Direction_UP,
	snake.DOWN:  websnake.Direction_DOWN,
	snake.LEFT:  websnake.Direction_LEFT,
	snake.RIGHT: websnake.Direction_RIGHT,
}

func playerToNetSnake(p game.Player) *websnake.GameState_Snake {
	netCoords := []*websnake.GameState_Coord{}
	netCoords = append(netCoords, &websnake.GameState_Coord{
		X: &p.Snake.Body[0].X,
		Y: &p.Snake.Body[0].Y,
	})

	for i, curr := range p.Snake.Body[1:] {
		x := curr.X - p.Snake.Body[i].X
		y := curr.Y - p.Snake.Body[i].Y
		netCoords = append(netCoords, &websnake.GameState_Coord{
			X: &x,
			Y: &y,
		})
	}

	webDir := DirToNetDir[p.Snake.Dir]

	log.Printf("Snake coords: %+v", netCoords)
	return &websnake.GameState_Snake{
		PlayerId: &p.Id,
		Points:   netCoords,
		// TODO: сейчас змея всегда живая, режим зомби не реализован
		State:         websnake.GameState_Snake_ALIVE.Enum(),
		HeadDirection: &webDir,
	}

}

var seq = int64(-1)

func generateSeq() int64 {
	seq++
	return seq
}

func (w *WebSnakeMasterNode) sendGameState(message pubsub.Message) {
	log.Println("Send gamestate----------------------------")
	defer log.Println("----------------------------Gamestate sent")

	netPlayers := []*websnake.GamePlayer{}
	for _, player := range w.game.Players {
		netPlayers = append(netPlayers, playerToNet(*player))
	}

	netGamePlayers := &websnake.GamePlayers{
		Players: netPlayers,
	}

	netSnakes := []*websnake.GameState_Snake{}
	for _, player := range w.game.Players {
		netSnakes = append(netSnakes, playerToNetSnake(*player))
	}

	netFood := []*websnake.GameState_Coord{}
	for i := range w.game.Food {
		netFood = append(netFood, &websnake.GameState_Coord{
			X: &w.game.Food[i].X,
			Y: &w.game.Food[i].Y,
		})
	}

	seq := generateSeq()
	msg := websnake.GameMessage{
		MsgSeq:     &seq,
		SenderId:   new(int32),
		ReceiverId: new(int32),
		Type: &websnake.GameMessage_State{
			State: &websnake.GameMessage_StateMsg{
				State: &websnake.GameState{
					Snakes:  netSnakes,
					Foods:   netFood,
					Players: netGamePlayers,
				},
			},
		},
	}
	data, err := proto.Marshal(&msg)
	if err != nil {
		log.Fatal(err)
	}

	// Рассылка по всем пользователям
	for _, player := range w.game.Players {
		to := net.UDPAddr{
			IP:   net.ParseIP(player.IpAddress),
			Port: int(player.Port),
			Zone: "",
		}
		// исключаем мастера из рассылки
		if player.Role != websnake.NodeRole_MASTER {
			go w.SendTo(data, &to)
		}
	}
}

func (w *WebSnakeMasterNode) DestroyNode() {
	w.conn.Close()
}

func playerToNet(player game.Player) *websnake.GamePlayer {
	return &websnake.GamePlayer{
		Name:      &player.Name,
		Id:        &player.Id,
		IpAddress: &player.IpAddress,
		Port:      &player.Port,
		Role:      &player.Role,
		Type:      &player.Type,
		Score:     &player.Score,
	}
}

func createAnnounce(game game.Game) websnake.GameAnnouncement {
	netPlayers := []*websnake.GamePlayer{}
	for _, p := range game.Players {
		netPlayers = append(netPlayers, playerToNet(*p))
	}

	canJoinBool := true
	return websnake.GameAnnouncement{
		Players: &websnake.GamePlayers{
			Players: netPlayers,
		},
		Config: &websnake.GameConfig{
			Width:        &game.GridSize.Width,
			Height:       &game.GridSize.Height,
			FoodStatic:   &game.StaticFood,
			StateDelayMs: &game.Delay,
		},
		CanJoin:  &canJoinBool,
		GameName: &game.Name,
	}
}

// Run in goroutine.
func (w *WebSnakeMasterNode) SendMultiAnnouncment() {
	log.Println("Send Announcement message")
	defer log.Println("---------------------------------------Announce message sent.")

	addr := "224.0.0.1"
	port := 8888

	for {
		announce := createAnnounce(*w.game)
		msg := websnake.GameMessage{
			MsgSeq:     new(int64),
			SenderId:   new(int32),
			ReceiverId: new(int32),
			Type: &websnake.GameMessage_Announcement{
				Announcement: &websnake.GameMessage_AnnouncementMsg{
					Games: []*websnake.GameAnnouncement{&announce},
				},
			},
		}
		data, err := proto.Marshal(&msg)
		if nil != err {
			log.Println("Catch err:", err)
			os.Exit(1)
		}
		updTimer := time.NewTimer(time.Second)
		<-updTimer.C
		w.SendMulticast(data, addr, port)
	}
}

func (w *WebSnakeMasterNode) SendTo(data []byte, to *net.UDPAddr) {
	log.Println("Send Unicast message to ", to.String())
	defer log.Println("---------------------------------------Unicast message sent.")

	if _, err := w.conn.WriteToUDP(data, to); err != nil {
		log.Println("Catch err: ", err)
		os.Exit(1)
	}
}

func (w *WebSnakeMasterNode) SendMulticast(data []byte, addr string, port int) {
	log.Println("Send Multicast message to ", w.conn.LocalAddr().Network())
	defer log.Println("---------------------------------------Multicast message sent.")

	// Send data to multicast group.
	netAddr := net.UDPAddr{
		IP:   net.ParseIP(addr),
		Port: port,
		Zone: "",
	}
	if _, err := w.conn.WriteTo(data, &netAddr); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (w *WebSnakeMasterNode) ListenAndServe() {
	log.Println("Master udp conn addr:", w.conn.LocalAddr().String())

	for {
		buf := make([]byte, 1024)
		n, addr, err := w.conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal("Err:", err)
		}
		buf = buf[0:n]

		message := websnake.GameMessage{}
		if err := proto.Unmarshal(buf, &message); err != nil {
			log.Fatal(err)
		}

		// log.Println("Read from udp socket in master:", message)

		eventMessage := pubsub.Message{
			Msg:  &message,
			From: addr,
			To:   nil,
		}
		switch {
		case message.GetAck() != nil:
			{
				// w.handleAck(eventMessage)
			}
		case message.GetJoin() != nil:
			{
				w.handleJoin(eventMessage)
			}
		case message.GetPing() != nil:
			{
				// HandlePing(eventMessage)
			}
		case message.GetSteer() != nil:
			{
				w.handleSteer(eventMessage)
			}
		default:
			log.Println("Handler for this type of message not implemented yet")
		}
	}
}

func (w *WebSnakeMasterNode) handleJoin(eventMessage pubsub.Message) {
	log.Println("Read from udp socket in master:", eventMessage)
	pubsub.GetGlobalPubSubService().Publish("join", eventMessage)
}

func (w *WebSnakeMasterNode) handleSteer(eventMessage pubsub.Message) {
	log.Println("Read from udp socket in master:", eventMessage)
	pubsub.GetGlobalPubSubService().Publish("steer", eventMessage)
}

func (w *WebSnakeMasterNode) sendAck(msg pubsub.Message) {
	log.Println("----------------------------send ack", msg.Msg)
	defer log.Println("----------------------sent ack")

	data, err := proto.Marshal(msg.Msg)
	if err != nil {
		log.Fatal(err)
	}

	w.SendTo(data, msg.To)
}
