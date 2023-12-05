package webnodes

import (
	"log"
	"main/websnake"

	"google.golang.org/protobuf/proto"
)

func (w *WebNode) buildAckBytes(receiverId int32) []byte {
	seq := generateSeq()
	msg := websnake.GameMessage{
		MsgSeq:     &seq,
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
	data, err := proto.Marshal(&msg)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
