package webnodes

import (
	"log"
	"main/game"
	"main/pubsub"
	"main/websnake"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/ipv4"
	"google.golang.org/protobuf/proto"
)

type WebNode struct {
	game      *game.Game
	conn      *net.UDPConn
	multiconn *mcastConn
}

func NewWebNode(game *game.Game) *WebNode {
	// Create unicast socket.
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create multicast socket.
	mcastgroup := "224.0.0.1"
	port := 8888
	mcastConn, err := createMulticastJoinedConnection(mcastgroup, port, "wlp2s0")
	if err != nil {
		log.Println("Catch error when create multicast packet connection:", err)
		os.Exit(1)
	}

	return &WebNode{
		conn:      conn,
		game:      game,
		multiconn: mcastConn,
	}

}

func (mcc *mcastConn) Close() {
	mcc.multiconn.Close()
}

func createMulticastJoinedConnection(mcastgroup string, port int, interfaceName string) (*mcastConn, error) {
	conn, err := net.ListenPacket("udp", mcastgroup+":"+strconv.Itoa(port))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	group := net.ParseIP(mcastgroup)
	if group == nil {
		log.Println("Invalid multicast group address.")
		return nil, err
	}

	p := ipv4.NewPacketConn(conn)
	netIface, err := net.InterfaceByName(interfaceName)
	log.Println(netIface)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	p.SetMulticastInterface(netIface)

	if err := p.JoinGroup(netIface, &net.UDPAddr{IP: group, Port: port}); err != nil {
		log.Println(err)
		return nil, err
	}

	return &mcastConn{
		multiconn: p,
		iface:     *netIface,
	}, nil
}

func (w *WebNode) RunLikeMaster() {
	go w.SendMultiAnnouncment()
	go w.ListenAndServe()
	go w.sendGameStates()
}

func (w *WebNode) RunLikeNormal() {
	go w.ReceiveMultiAnnouncments()
	go w.ListenAndServe()
}

func (w *WebNode) mapModelPlayersToNetPlayers() *websnake.GamePlayers {
	netPlayers := []*websnake.GamePlayer{}
	for _, player := range w.game.Players {
		netPlayers = append(netPlayers, playerToNet(*player))
	}

	return &websnake.GamePlayers{
		Players: netPlayers,
	}
}

func (w *WebNode) mapModelPlayersToNetSnakes() []*websnake.GameState_Snake {
	netSnakes := []*websnake.GameState_Snake{}
	for _, player := range w.game.Players {
		if player.Snake.IsAlive {
			netSnakes = append(netSnakes, playerToNetSnake(*player))
		}
	}

	return netSnakes
}

func (w *WebNode) mapModelFoodToNetFood() []*websnake.GameState_Coord {
	netFood := []*websnake.GameState_Coord{}
	for i := range w.game.Food {
		netFood = append(netFood, &websnake.GameState_Coord{
			X: &w.game.Food[i].X,
			Y: &w.game.Food[i].Y,
		})
	}

	return netFood
}

func (w *WebNode) sendGameStates() {
	for {
		<-time.NewTimer(time.Duration(w.game.Delay) * time.Millisecond).C
		w.sendGameState()
	}
}

func (w *WebNode) sendGameState() {
	if len(w.game.Players) < 2 {
		return
	}

	globalStateOrder++

	log.Println("Send gamestate----------------------------")
	defer log.Println("----------------------------Gamestate sent")

	// Рассылка по всем пользователям
	for _, player := range w.game.Players {
		to := net.UDPAddr{
			IP:   net.ParseIP(player.IpAddress),
			Port: int(player.Port),
			Zone: "",
		}

		// исключаем мастера из рассылки
		if player.Role != websnake.NodeRole_MASTER {
			go w.SendTo(w.buildGameStateBytes(player.Id, globalStateOrder), &to)
		}
	}
}

// Run in goroutine.
func (w *WebNode) ReceiveMultiAnnouncments() {
	for {
		buf := make([]byte, 1024)
		n, _, srcAddr, err := w.multiconn.multiconn.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		buf = buf[0:n]

		// announce := websnake.GameAnnouncement{}
		msg := websnake.GameMessage{}
		if err := proto.Unmarshal(buf, &msg); err != nil {
			log.Println(err)
			os.Exit(1)
		}

		// Определение типа сообщения
		an := msg.GetAnnouncement()
		addMasterAddresIntoPlayers(an.GetGames()[0].Players.Players, srcAddr)

		addrAndPort := strings.Split(srcAddr.String(), ":")
		addr := addrAndPort[0]
		port, err := strconv.Atoi(addrAndPort[1])
		if err != nil {
			log.Fatal(err)
		}

		// тут происходит публикование сообщения о текущих играх по внутренней шине сообщений.
		pubsub.GetGlobalPubSubService().Publish("announce", pubsub.Message{
			Msg: &msg,
			From: &net.UDPAddr{
				IP:   net.ParseIP(addr),
				Port: port,
				Zone: "",
			},
			To: nil,
		})
	}
}

func (w *WebNode) ListenAndServe() {
	for {
		buf := make([]byte, 1024*8)
		n, addr, err := w.conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal("Err:", err)
		}
		buf = buf[0:n]

		message := websnake.GameMessage{}
		if err := proto.Unmarshal(buf, &message); err != nil {
			log.Fatal(err)
		}

		eventMessage := pubsub.Message{
			Msg:  &message,
			From: addr,
			To:   nil,
		}
		switch {
		case message.GetAck() != nil:
			{
				w.handleAck(eventMessage)
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
		case message.GetState() != nil:
			{
				w.handleState(eventMessage)
			}
		default:
			log.Println("Handler for this type of message not implemented yet")
		}
	}

}

func (w *WebNode) handleAck(message pubsub.Message) {
	switch w.game.GetMainPlayer().Role {
	case websnake.NodeRole_MASTER:
		{
			w.game.UpdateUserTimeout(message.Msg.GetSenderId(), time.Now())
		}
	default:
		{
			// Любой узел, кроме мастера, получает ack в двух случаях:
			// 1. Подтверждение соединения
			// 2. Подтверждение любого другого сообщения.
			if !w.game.IsRun {
				w.game.SetMainPlayer(*message.Msg.ReceiverId)
				w.game.Run()
			}
		}
	}

	// w.messageQueue.Dequeue(message.Msg.GetMsgSeq())
}

func (w *WebNode) handleJoin(message pubsub.Message) {
	switch w.game.GetMainPlayer().Role {
	case websnake.NodeRole_MASTER:
		{
			join := message.Msg.GetJoin()
			senderIP := message.From.IP
			senderPort := message.From.Port
			newPlayer, err := w.game.AddPlayer(join.GetGameName(), senderIP.String(), senderPort, websnake.NodeRole_NORMAL, websnake.PlayerType_HUMAN)
			if err != nil {
				// TODO: build and send error
			}

			w.SendTo(w.buildAckBytes(newPlayer.Id), message.From)

			// w.game.UpdateUserTimeout(message.Msg.GetSenderId(), time.Now())
		}
	default:
		{
		}
	}

	// w.messageQueue.Dequeue(message.Msg.GetMsgSeq())
}

func (w *WebNode) handleSteer(message pubsub.Message) {
	switch w.game.GetMainPlayer().Role {
	case websnake.NodeRole_MASTER:
		{
			log.Println("Read from udp socket in master:", message)
			pubsub.GetGlobalPubSubService().Publish("steer", message)
			// w.game.UpdateUserTimeout(message.Msg.GetSenderId(), time.Now())
		}
	default:
		{
		}
	}

	// w.messageQueue.Dequeue(message.Msg.GetMsgSeq())
}

func (w *WebNode) handleState(message pubsub.Message) {
	switch w.game.GetMainPlayer().Role {
	case websnake.NodeRole_MASTER:
		{
		}
	default:
		{
			log.Println("Read from udp socket:", message)
			players := message.Msg.GetState().State.Players.Players
			addMasterAddresIntoPlayers(players, message.From)
			pubsub.GetGlobalPubSubService().Publish("newgamestate", message)
		}
	}

	// w.messageQueue.Dequeue(message.Msg.GetMsgSeq())
}

// Run in goroutine.
func (w *WebNode) SendMultiAnnouncment() {
	log.Println("Send Announcement message")
	defer log.Println("---------------------------------------Announce message sent.")

	addr := "224.0.0.1"
	port := 8888

	to := &net.UDPAddr{
		IP:   net.ParseIP(addr),
		Port: port,
		Zone: "",
	}

	for {
		announce := createAnnounce(*w.game)
		seq := generateSeq()
		msg := websnake.GameMessage{
			MsgSeq:     &seq,
			SenderId:   w.game.MainPlayerID,
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
		w.SendTo(data, to)
	}
}

func (w *WebNode) SendTo(data []byte, to *net.UDPAddr) {
	log.Println("Send Unicast message to ", to.String())
	defer log.Println("---------------------------------------Unicast message sent.")

	if _, err := w.conn.WriteToUDP(data, to); err != nil {
		log.Println("Catch err: ", err)
		os.Exit(1)
	}
}
