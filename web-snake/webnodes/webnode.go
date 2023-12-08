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

var seq = int64(-1)

func generateSeq() int64 {
	seq++
	return seq
}

func (w *WebNode) SetGame(g *game.Game) {
	w.game = g
}

func (w *WebNode) Join() {
	announce := <-w.game.ConnectionChannel

	playerType := websnake.PlayerType_HUMAN
	playerRole := websnake.NodeRole_NORMAL

	name := w.game.GetMainPlayer().Name
	gameName := announce.GetGameName()

	w.SendTo(w.buildJoinBytes(name, gameName, playerType, playerRole), getMasterAddressFromAnnounce(announce))
}

func FindPlayerWithRole(g *game.Game, role websnake.NodeRole) *game.Player {
	for _, p := range g.Players {
		log.Println(*p)
		if p.Role == role {
			return p
		}
	}

	return nil
}

func (w *WebNode) GetMasterAddress() *net.UDPAddr {
	for _, p := range w.game.Players {
		log.Println("try to find master:", *p)
	}

	master := FindPlayerWithRole(w.game, websnake.NodeRole_MASTER)
	log.Println("Master:", master)
	if master != nil {
		return &net.UDPAddr{
			IP:   net.ParseIP(master.IpAddress),
			Port: int(master.Port),
			Zone: "",
		}
	}

	deputy := FindPlayerWithRole(w.game, websnake.NodeRole_DEPUTY)
	if deputy != nil {
		return &net.UDPAddr{
			IP:   net.ParseIP(deputy.IpAddress),
			Port: int(deputy.Port),
			Zone: "",
		}
	}
	log.Println("Deputy:", deputy)

	return nil
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

type mcastConn struct {
	multiconn *ipv4.PacketConn
	iface     net.Interface
}

func NewEmptyWebNode() *WebNode {
	// Create unicast socket.
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create multicast socket.
	mcastgroup := "239.192.0.4"
	port := 9192
	mcastConn, err := createMulticastJoinedConnection(mcastgroup, port, "wlp2s0")
	if err != nil {
		log.Println("Catch error when create multicast packet connection:", err)
		os.Exit(1)
	}
	node := &WebNode{
		conn:      conn,
		multiconn: mcastConn,
	}

	steerSender := pubsub.Subscriber{
		EventChannel: make(chan string),
		EventHandler: func(msg pubsub.Message) {
			node.sendSteer(msg)
		},
	}

	pubsub.GetGlobalPubSubService().Subscribe("steersend", steerSender)

	return node
}

func NewWebNode(game *game.Game) *WebNode {
	// Create unicast socket.
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create multicast socket.
	mcastgroup := "239.192.0.4"
	port := 9192
	mcastConn, err := createMulticastJoinedConnection(mcastgroup, port, "wlp2s0")
	if err != nil {
		log.Println("Catch error when create multicast packet connection:", err)
		os.Exit(1)
	}
	node := &WebNode{
		conn:      conn,
		game:      game,
		multiconn: mcastConn,
	}

	steerSender := pubsub.Subscriber{
		EventChannel: make(chan string),
		EventHandler: func(msg pubsub.Message) {
			node.sendSteer(msg)
		},
	}

	pubsub.GetGlobalPubSubService().Subscribe("steersend", steerSender)

	return node
}

func (w *WebNode) sendSteer(msg pubsub.Message) {
	log.Println("---------------------------------Steer send to", msg.To)
	w.SendTo(w.buildSteerBytes(msg.Msg.GetSteer().GetDirection()), w.GetMasterAddress())
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
	go w.CleanupPlayers()
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
		// if player.Snake.IsAlive {
		netSnakes = append(netSnakes, playerToNetSnake(*player))
		// }
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
		log.Println("send game state -----------------------------------")
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

func addMasterAddresIntoPlayers(players []*websnake.GamePlayer, srcAddr net.Addr) []*websnake.GamePlayer {
	for i, player := range players {
		if *player.Role == *websnake.NodeRole_MASTER.Enum() {
			addrPort := strings.Split(srcAddr.String(), ":")
			players[i].IpAddress = &addrPort[0]
			port, err := strconv.Atoi(addrPort[1])
			port32 := int32(port)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

			players[i].Port = &port32
		}
	}

	return players
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

func IsPlayerTimeIsOut(player *game.Player, deleyInMs int32) bool {
	return time.Since(player.LastTimeout).Milliseconds() > int64(deleyInMs)
}

func (w *WebNode) CleanupPlayers() {
	for {
		<-time.NewTimer(time.Second).C
		switch w.game.GetMainPlayer().Role {
		case websnake.NodeRole_MASTER:
			{
				for _, player := range w.game.Players {
					if player.Role != websnake.NodeRole_MASTER {
						// теряем бойца.
						if IsPlayerTimeIsOut(player, w.game.Delay*4) {
							if player.Role == websnake.NodeRole_DEPUTY {
								// удаляем депути
								delete(w.game.Players, player.Id)
								// Ставим нового
								w.setDeputy()
							}

							player.Snake.IsZombie = true
						}
					}
				}
			}

		case websnake.NodeRole_NORMAL:
			{
				master := FindPlayerWithRole(w.game, websnake.NodeRole_MASTER)
				if master != nil && IsPlayerTimeIsOut(master, w.game.Delay*4) {
					// удаляем мастера.
					delete(w.game.Players, master.Id)
				}
			}
		case websnake.NodeRole_DEPUTY:
			{
				master := FindPlayerWithRole(w.game, websnake.NodeRole_MASTER)
				if master != nil && IsPlayerTimeIsOut(master, w.game.Delay*4) {
					log.Println("delete master by time out")
					// удаляем мастера.
					master.Snake.IsZombie = true
					// delete(w.game.Players, master.Id)
					w.game.Players[*w.game.MainPlayerID].Role = websnake.NodeRole_MASTER
					// start game like master
					go w.game.MainLoop()
					go w.SendMultiAnnouncment()
					go w.sendGameStates()
				}
			}
		}
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

		// Игнорируем зомби челов.
		if senderPlayer, ok := w.game.Players[message.GetSenderId()]; ok {
			if senderPlayer.Snake.IsZombie {
				continue
			}
		}

		switch {
		case message.GetAck() != nil:
			{
				log.Println("Catch ACK-------------------------------")
				w.handleAck(eventMessage)
			}
		case message.GetJoin() != nil:
			{
				log.Println("Catch JOIN-------------------------------")
				w.handleJoin(eventMessage)
			}
		case message.GetPing() != nil:
			{
				w.handlePing(eventMessage)
			}
		case message.GetSteer() != nil:
			{
				w.handleSteer(eventMessage)
			}
		case message.GetState() != nil:
			{
				log.Println("Catch GAMESTATE-------------------------------")
				w.handleState(eventMessage)
			}
		default:
			log.Println("Handler for this type of message not implemented yet")
		}
	}

}

var wasJoined = false

func (w *WebNode) handleAck(message pubsub.Message) {
	w.game.UpdateUserTimeout(message.Msg.GetSenderId(), time.Now())
	switch w.game.GetMainPlayer().Role {
	case websnake.NodeRole_MASTER:
		{
		}
	default:
		{
			// Любой узел, кроме мастера, получает ack в двух случаях:
			// 1. Подтверждение соединения.
			// 2. Подтверждение любого другого сообщения.
			// if !w.game.IsRun {
			if !wasJoined {
				w.game.SetMainPlayer(*message.Msg.ReceiverId)
				w.game.Run()
				wasJoined = true
			}
			// }
		}
	}

	// w.messageQueue.Dequeue(message.Msg.GetMsgSeq())
}

func (w *WebNode) handlePing(message pubsub.Message) {
	switch w.game.GetMainPlayer().Role {
	case websnake.NodeRole_MASTER:
		{
			w.game.UpdateUserTimeout(message.Msg.GetSenderId(), time.Now())
			// send ack
			w.SendTo(w.buildAckBytes(message.Msg.GetMsgSeq(), message.Msg.GetSenderId()), message.From)
		}
	default:
		{
		}
	}

	// w.messageQueue.Dequeue(message.Msg.GetMsgSeq())
}

func (w *WebNode) setDeputy() *game.Player {
	currentDeputy := FindPlayerWithRole(w.game, websnake.NodeRole_DEPUTY)
	if currentDeputy != nil {
		return currentDeputy
	}

	for _, player := range w.game.Players {
		if player.Role == websnake.NodeRole_NORMAL {
			player.Role = websnake.NodeRole_DEPUTY
			return player
		}
	}

	return nil
}

func (w *WebNode) handleJoin(message pubsub.Message) {
	switch w.game.GetMainPlayer().Role {
	case websnake.NodeRole_MASTER:
		{
			join := message.Msg.GetJoin()
			senderIP := message.From.IP
			senderPort := message.From.Port
			log.Println("New player name from join message:", join.GetGameName())

			newPlayer := &game.Player{}
			var err error
			if FindPlayerWithRole(w.game, websnake.NodeRole_DEPUTY) == nil {
				newPlayer, err = w.game.AddPlayer(join.GetPlayerName(), senderIP.String(), senderPort, websnake.NodeRole_DEPUTY, websnake.PlayerType_HUMAN)
			} else {
				newPlayer, err = w.game.AddPlayer(join.GetPlayerName(), senderIP.String(), senderPort, websnake.NodeRole_NORMAL, websnake.PlayerType_HUMAN)
			}

			if err != nil {
				// TODO: build and send error
			}

			// send ack
			w.SendTo(w.buildAckBytes(message.Msg.GetMsgSeq(), newPlayer.Id), message.From)

			w.game.UpdateUserTimeout(message.Msg.GetSenderId(), time.Now())
		}
	default:
		{
		}
	}

	// w.messageQueue.Dequeue(message.Msg.GetMsgSeq())
}

func (w *WebNode) handleSteer(message pubsub.Message) {
	w.game.UpdateUserTimeout(message.Msg.GetSenderId(), time.Now())
	switch w.game.GetMainPlayer().Role {
	case websnake.NodeRole_MASTER:
		{
			pubsub.GetGlobalPubSubService().Publish("steer", message)
			playerId := message.Msg.SenderId
			log.Println("----------------------------Steer player with id:", *playerId)
			steer := message.Msg.GetSteer()
			w.game.SteerPlayerSnake(*playerId, game.NetDirToModel[steer.GetDirection()])

			// send ack
			w.SendTo(w.buildAckBytes(message.Msg.GetMsgSeq(), message.Msg.GetSenderId()), message.From)
		}
	default:
		{
		}
	}

	// w.messageQueue.Dequeue(message.Msg.GetMsgSeq())
}

func (w *WebNode) handleState(message pubsub.Message) {
	w.game.UpdateUserTimeout(message.Msg.GetSenderId(), time.Now())
	switch w.game.GetMainPlayer().Role {
	case websnake.NodeRole_MASTER:
		{
		}
	default:
		{
			players := message.Msg.GetState().State.Players.Players
			addMasterAddresIntoPlayers(players, message.From)
			for _, p := range w.game.Players {
				log.Println("game player:", p)
			}
			w.game.UpdateGameState(message)
			// send ack
			w.SendTo(w.buildAckBytes(message.Msg.GetMsgSeq(), message.Msg.GetSenderId()), message.From)

			// w.game.UpdateGameState(message)
			for i := 0; i < 2; i++ {
				<-time.NewTimer(time.Duration(w.game.Delay / 2)).C
				w.SendTo(w.buildPingBytes(), message.From)
			}
		}
	}
	// w.messageQueue.Dequeue(message.Msg.GetMsgSeq())
}

// Run in goroutine.
func (w *WebNode) SendMultiAnnouncment() {
	log.Println("Send Announcement message")
	defer log.Println("---------------------------------------Announce message sent.")

	addr := "239.192.0.4"
	port := 9192

	to := &net.UDPAddr{
		IP:   net.ParseIP(addr),
		Port: port,
		Zone: "",
	}

	for {
		<-time.NewTimer(time.Second).C
		w.SendTo(w.buildAnnounceBytes(w.createAnnounce()), to)
	}
}

func (w *WebNode) createAnnounce() websnake.GameAnnouncement {
	netPlayers := []*websnake.GamePlayer{}
	for _, p := range w.game.Players {
		netPlayers = append(netPlayers, playerToNet(*p))
	}

	canJoinBool := true
	return websnake.GameAnnouncement{
		Players: &websnake.GamePlayers{
			Players: netPlayers,
		},
		Config: &websnake.GameConfig{
			Width:        &w.game.GridSize.Width,
			Height:       &w.game.GridSize.Height,
			FoodStatic:   &w.game.StaticFood,
			StateDelayMs: &w.game.Delay,
		},
		CanJoin:  &canJoinBool,
		GameName: &w.game.Name,
	}
}

func (w *WebNode) SendTo(data []byte, to *net.UDPAddr) {
	log.Println("Send Unicast message to ", to.String())
	defer log.Println("---------------------------------------Unicast message sent.")

	if _, err := w.conn.WriteToUDP(data, to); err != nil {
		log.Fatal("Catch err: ", err)
	}
}
