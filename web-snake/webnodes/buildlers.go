package webnodes

import (
	"log"
	"main/websnake"

	"google.golang.org/protobuf/proto"
)

func (w *WebNode) buildAckBytes(senderSeqId int64, receiverId int32) []byte {
	log.Println("Send Ack---------------------------------", senderSeqId, receiverId)
	msg := websnake.GameMessage{
		MsgSeq:     &senderSeqId,
		SenderId:   &w.game.GetMainPlayer().Id,
		ReceiverId: &receiverId,
		Type: &websnake.GameMessage_Ack{
			Ack: &websnake.GameMessage_AckMsg{},
		},
	}

	data, err := proto.Marshal(&msg)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func (w *WebNode) buildAnnounceBytes(announce websnake.GameAnnouncement) []byte {
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
		log.Fatal(err)
	}

	return data
}

func (w *WebNode) buildJoinBytes(username string, gamename string, playerType websnake.PlayerType, playerRole websnake.NodeRole) []byte {
	seq := generateSeq()
	message := websnake.GameMessage{
		MsgSeq: &seq,
		Type: &websnake.GameMessage_Join{
			Join: &websnake.GameMessage_JoinMsg{
				PlayerType:    &playerType,
				PlayerName:    &username,
				GameName:      &gamename,
				RequestedRole: &playerRole,
			},
		},
	}

	data, err := proto.Marshal(&message)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func (w *WebNode) buildPingBytes() []byte {
	seq := generateSeq()
	gameMessage := websnake.GameMessage{
		MsgSeq:     &seq,
		SenderId:   w.game.MainPlayerID,
		ReceiverId: new(int32),
		Type:       &websnake.GameMessage_Ping{Ping: &websnake.GameMessage_PingMsg{}},
	}

	data, err := proto.Marshal(&gameMessage)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

var globalStateOrder = int32(0)

func (w *WebNode) buildSteerBytes(direction websnake.Direction) []byte {
	seq := generateSeq()
	log.Println("In steer builder main player id:", *w.game.MainPlayerID)

	gameMessage := websnake.GameMessage{
		MsgSeq:     &seq,
		SenderId:   w.game.MainPlayerID,
		ReceiverId: new(int32),
		Type: &websnake.GameMessage_Steer{
			Steer: &websnake.GameMessage_SteerMsg{
				Direction: &direction,
			},
		},
	}

	data, err := proto.Marshal(&gameMessage)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func (w *WebNode) buildGameStateBytes(receiverId int32, stateOrder int32) []byte {
	seq := generateSeq()
	msg := websnake.GameMessage{
		MsgSeq:     &seq,
		SenderId:   w.game.MainPlayerID,
		ReceiverId: &receiverId,
		Type: &websnake.GameMessage_State{
			State: &websnake.GameMessage_StateMsg{
				State: &websnake.GameState{
					StateOrder: &stateOrder,
					Snakes:     w.mapModelPlayersToNetSnakes(),
					Foods:      w.mapModelFoodToNetFood(),
					Players:    w.mapModelPlayersToNetPlayers(),
				},
			},
		},
	}
	log.Println("players in gamestate message:", msg.GetState().State.GetPlayers())
	data, err := proto.Marshal(&msg)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
