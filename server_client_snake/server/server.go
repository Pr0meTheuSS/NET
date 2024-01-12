package server

import (
	"errors"
	"log"
	"net"
	"snake/game"
	"snake/models"
	"snake/websnake_proto_gen/main/websnake"

	"google.golang.org/protobuf/proto"
)

const (
	maxDatagramSize = 65507
)

type IServer interface {
	ListenAndServe() error
}

type Server struct {
	port           int
	connection     net.UDPConn
	gameController game.IGameController
	playersToGames map[string]*game.Game
}

func NewServer(port int) *Server {
	return &Server{
		port:           port,
		connection:     net.UDPConn{},
		gameController: &game.GameControllerImpl{},
		playersToGames: map[string]*game.Game{},
	}
}

func (s *Server) ListenAndServe() error {
	log.Println("---------------Server listening on port:", s.port)

	go s.sendQueuedMessages()

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: s.port,
	})

	if err != nil {
		return err
	}
	defer conn.Close()

	s.connection = *conn

	for {
		buffer := make([]byte, maxDatagramSize)
		n, from, err := conn.ReadFromUDP(buffer)

		if err != nil {
			log.Println(err.Error())
			continue
		}
		buffer = buffer[:n]
		log.Printf("Server recv datagram from %+v\n", from)

		if game, ok := s.playersToGames[from.String()]; ok {
			s.handleGameMessage(game, buffer, from)
		} else {
			s.handleConnection(buffer, from)
		}
	}
}

func (s *Server) sendQueuedMessages() {
	for {
		for _, g := range s.gameController.GetAllGames() {
			go func(game *game.Game) {
				msg, status := <-game.MessageChannel
				if status {
					data, err := proto.Marshal(msg.GameMessage)
					if err != nil {
						log.Println()
					}

					s.sendTo(data, msg.To)
				}
			}(g)
		}
	}
}

func (s *Server) handleGameMessage(game *game.Game, datagram []byte, from *net.UDPAddr) {
	message := websnake.GameMessage{}
	if err := proto.Unmarshal(datagram, &message); err != nil {
		log.Println(err.Error())
		return
	}

	switch {
	case message.GetAck() != nil:
		{
			s.handleAck(&message, from, game)
		}
	case message.GetPing() != nil:
		{
			s.handlePing(&message, from, game)
		}
	case message.GetSteer() != nil:
		{
			s.handleSteer(&message, from, game)
		}
	// TODO: implement.
	// Not supported yet.
	// case message.GetRoleChange() != nil:
	// 	{
	// 		handleRoleChange(&message, from, game)
	// 	}
	// TODO: implement.
	// Not supported yet.
	// case message.GetState() != nil:
	// 	{
	// 		handleState(&message, from, game)
	// 	}
	default:
		log.Println("Handler for this type of message not implemented yet")
	}
}

func (s *Server) handleAck(ack *websnake.GameMessage, from *net.UDPAddr, game *game.Game) {
	game.UpdatePlayerByAddress(from)

	// TODO: implement.
	// Not supported yet.
	// messageQueue.Dequeue(ack.GetMsgSeq())
}

func (s *Server) handlePing(ping *websnake.GameMessage, from *net.UDPAddr, game *game.Game) {
	game.UpdatePlayerByAddress(from)
	s.sendAckToMessage(ping, from)
}

func (s *Server) handleSteer(steer *websnake.GameMessage, from *net.UDPAddr, game *game.Game) {
	game.UpdatePlayerByAddress(from)
	s.sendAckToMessage(steer, from)
	netDirToModelDir := map[websnake.Direction]models.Direction{
		websnake.Direction_DOWN:  models.DOWN,
		websnake.Direction_UP:    models.UP,
		websnake.Direction_LEFT:  models.LEFT,
		websnake.Direction_RIGHT: models.RIGHT,
	}

	game.SteerPlayerByAddress(from, netDirToModelDir[steer.GetSteer().GetDirection()])
}

func (s *Server) sendAckToMessage(msg *websnake.GameMessage, to *net.UDPAddr) {
	ackGameMessage := websnake.GameMessage{
		MsgSeq:     msg.MsgSeq,
		SenderId:   new(int32),
		ReceiverId: new(int32),
		Type:       &websnake.GameMessage_Ack{Ack: &websnake.GameMessage_AckMsg{}},
	}

	data, err := proto.Marshal(&ackGameMessage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	s.sendTo(data, to)
}

func (s *Server) handleConnection(datagram []byte, from *net.UDPAddr) {
	log.Println("---------------Server handle new connection")
	defer log.Println("Server handled new connection---------------")

	gameMessage := websnake.GameMessage{}
	if err := proto.Unmarshal(datagram, &gameMessage); err != nil {
		log.Println(err.Error())
		return
	}

	switch {
	case gameMessage.GetJoin() != nil:
		{
			s.handleJoinGameMessage(gameMessage.GetJoin(), from)
		}
	case gameMessage.GetDiscover() != nil:
		{
			s.handleDiscoverGameMessage(gameMessage.GetDiscover(), from)
		}
	default:
		{
			log.Println(errors.New("cannot handle message: wrong type of message"))
		}
	}
}

func (s *Server) handleJoinGameMessage(joinMessage *websnake.GameMessage_JoinMsg, from *net.UDPAddr) {
	// TODO: add error checking
	gameToJoin, _ := s.gameController.FindGameByName(joinMessage.GetGameName())

	if gameToJoin == nil {
		// TODO: add error checking
		newGame := s.gameController.CreateNewGame(joinMessage.GetGameName(), game.GenerateRandomConfig())
		s.playersToGames[from.String()] = newGame

		newGame.JoinPlayer(joinMessage.GetPlayerName(), joinMessage.GetPlayerType(), joinMessage.GetRequestedRole(), from)
		newGame.Run()
	} else {
		// TODO: add error checking
		gameToJoin.JoinPlayer(joinMessage.GetPlayerName(), joinMessage.GetPlayerType(), joinMessage.GetRequestedRole(), from)
		s.playersToGames[from.String()] = gameToJoin
	}

	ackGameMessage := websnake.GameMessage{
		MsgSeq:     new(int64),
		SenderId:   new(int32),
		ReceiverId: new(int32),
		Type:       &websnake.GameMessage_Ack{Ack: &websnake.GameMessage_AckMsg{}},
	}

	data, err := proto.Marshal(&ackGameMessage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	s.sendTo(data, from)
}

func (s *Server) handleDiscoverGameMessage(discoverMessage *websnake.GameMessage_DiscoverMsg, from *net.UDPAddr) {
	games := s.gameController.GetAllGames()
	announceMsg := buildAnnounceMsgByGames(games)

	gameMessage := websnake.GameMessage{
		MsgSeq:     new(int64),
		SenderId:   new(int32),
		ReceiverId: new(int32),
		Type: &websnake.GameMessage_Announcement{
			Announcement: announceMsg,
		},
	}

	data, err := proto.Marshal(&gameMessage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	s.sendTo(data, from)
}

func (s *Server) sendTo(data []byte, to *net.UDPAddr) {
	log.Println("Send Unicast message to ", to.String())
	defer log.Println("---------------------------------------Unicast message sent.")

	if _, err := s.connection.WriteToUDP(data, to); err != nil {
		log.Fatal("Catch err: ", err)
	}
}

func buildAnnounceMsgByGames(games []*game.Game) *websnake.GameMessage_AnnouncementMsg {
	announce := websnake.GameMessage_AnnouncementMsg{
		Games: []*websnake.GameAnnouncement{},
	}

	for _, g := range games {
		announce.Games = append(announce.Games, buildGameAnnouncement(g))
	}

	return &announce
}

func buildGameAnnouncement(g *game.Game) *websnake.GameAnnouncement {
	gameName := g.Name
	canJoin := true
	gamePlayers := playersToGamePlayers(g.Players)
	gameConfig := gameToGameConfig(g)

	return &websnake.GameAnnouncement{
		Players:  gamePlayers,
		Config:   &gameConfig,
		CanJoin:  &canJoin,
		GameName: &gameName,
	}
}

func gameToGameConfig(g *game.Game) websnake.GameConfig {
	width := g.GridSize.Width
	height := g.GridSize.Height
	foodStatic := g.StaticFood
	stateDelayMs := g.Delay

	return websnake.GameConfig{
		Width:        &width,
		Height:       &height,
		FoodStatic:   &foodStatic,
		StateDelayMs: &stateDelayMs,
	}
}

func playersToGamePlayers(players map[int32]*game.Player) *websnake.GamePlayers {
	netPlayers := []*websnake.GamePlayer{}
	for _, player := range players {
		netPlayers = append(netPlayers, playerToNet(*player))
	}

	return &websnake.GamePlayers{
		Players: netPlayers,
	}
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
