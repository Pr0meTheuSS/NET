// server_test.go
package server

import (
	"net"
	"snake/websnake_proto_gen/main/websnake"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"
)

func init() {
	s := NewServer(9192)

	go func() {
		s.ListenAndServe()
	}()

	time.Sleep(time.Millisecond * 100)
}

func TestUDPServerConnection(t *testing.T) {
	conn, err := net.Dial("udp", "127.0.0.1:9192")
	if err != nil {
		t.Fatalf("Error connecting to UDP server: %v", err)
	}
	defer conn.Close()
}

func TestUDPServerSendDiscoverMessage(t *testing.T) {
	conn, err := net.Dial("udp", "127.0.0.1:9192")
	if err != nil {
		t.Fatalf("Error connecting to UDP server: %v", err)
	}
	defer conn.Close()

	message := websnake.GameMessage{
		MsgSeq:     new(int64),
		SenderId:   new(int32),
		ReceiverId: new(int32),
		Type:       &websnake.GameMessage_Discover{Discover: &websnake.GameMessage_DiscoverMsg{}},
	}

	data, err := proto.Marshal(&message)
	if err != nil {
		t.Fatalf("Error marshal test discover message\n")
	}

	n, err := conn.Write(data)
	if n != len(data) || err != nil {
		t.Fatalf("Wrong amount of sent bytes or error while sending")
	}
}

func TestUDPServerSendDiscoverMessageAndRecevieGameAnnounce(t *testing.T) {
	conn, err := net.Dial("udp", "127.0.0.1:9192")
	if err != nil {
		t.Fatalf("Error connecting to UDP server: %v", err)
	}
	defer conn.Close()

	message := websnake.GameMessage{
		MsgSeq:     new(int64),
		SenderId:   new(int32),
		ReceiverId: new(int32),
		Type:       &websnake.GameMessage_Discover{Discover: &websnake.GameMessage_DiscoverMsg{}},
	}

	data, err := proto.Marshal(&message)
	if err != nil {
		t.Fatalf("Error marshal test discover message\n")
	}

	n, err := conn.Write(data)
	if n != len(data) || err != nil {
		t.Fatalf("Wrong amount of sent bytes or error while sending")
	}

	buffer := make([]byte, 1024*64)
	n, err = conn.Read(buffer)
	if err != nil {
		t.Fatalf(err.Error())
	}

	buffer = buffer[:n]

	gameMessage := websnake.GameMessage{}
	if err = proto.Unmarshal(buffer, &gameMessage); err != nil {
		t.Fatalf(err.Error())
	}

	if gameMessage.GetAnnouncement() == nil {
		t.Fatalf("Expected game message with announce of games, but received something else")
	}
}

func TestUDPServerSendJoinMessageAndReceiveAckOrError(t *testing.T) {
	conn, err := net.Dial("udp", "127.0.0.1:9192")
	if err != nil {
		t.Fatalf("Error connecting to UDP server: %v", err)
	}
	defer conn.Close()

	message := websnake.GameMessage{
		MsgSeq:     new(int64),
		SenderId:   new(int32),
		ReceiverId: new(int32),
		Type: &websnake.GameMessage_Join{Join: &websnake.GameMessage_JoinMsg{
			PlayerType:    websnake.PlayerType_HUMAN.Enum(),
			PlayerName:    new(string),
			GameName:      new(string),
			RequestedRole: websnake.NodeRole_NORMAL.Enum(),
		}},
	}

	data, err := proto.Marshal(&message)
	if err != nil {
		t.Fatalf("Error marshal test join message\n")
	}

	n, err := conn.Write(data)
	if n != len(data) || err != nil {
		t.Fatalf("Wrong amount of sent bytes or error while sending")
	}

	buffer := make([]byte, 1024*64)
	n, err = conn.Read(buffer)
	if err != nil {
		t.Fatalf(err.Error())
	}

	buffer = buffer[:n]

	gameMessage := websnake.GameMessage{}
	if err = proto.Unmarshal(buffer, &gameMessage); err != nil {
		t.Fatalf(err.Error())
	}

	if (gameMessage.GetAck() != nil) == (gameMessage.GetError() != nil) {
		t.Fatalf("Expected game message with ack or error message, but received something else")
	}
}
