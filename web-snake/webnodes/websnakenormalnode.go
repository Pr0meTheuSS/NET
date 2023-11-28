package webnodes

import (
	"log"
	"main/pubsub"
	"main/websnake"
	"net"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/ipv4"
	"google.golang.org/protobuf/proto"
)

type mcastConn struct {
	baseconn  net.PacketConn
	multiconn *ipv4.PacketConn
	iface     net.Interface
}

type webSnakeNormalNode struct {
	username    string
	mcastconn   *mcastConn
	unicastconn *net.UDPConn
}

func NewWebSnakeNormalNode(username string) *webSnakeNormalNode {
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Local address:", conn.LocalAddr())

	// log.Println("Recv Announcement message")
	mcastgroup := "239.192.0.4"
	port := 9192
	mcastConn, err := createMulticastJoinedConnection(mcastgroup, port, "wlp2s0")
	if err != nil {
		log.Println("Catch error when create multicast packet connection:", err)
		os.Exit(1)
	}
	node := &webSnakeNormalNode{
		username:    username,
		mcastconn:   mcastConn,
		unicastconn: conn,
	}

	joiner := pubsub.Subscriber{
		EventChannel: make(chan string),
		EventHandler: func(msg pubsub.Message) {
			node.join(msg)
		},
	}

	pubsub.GetGlobalPubSubService().Subscribe("connection", joiner)

	steerSender := pubsub.Subscriber{
		EventChannel: make(chan string),
		EventHandler: func(msg pubsub.Message) {
			node.sendSteer(msg)
		},
	}

	pubsub.GetGlobalPubSubService().Subscribe("steersend", steerSender)

	return node
}

func (w *webSnakeNormalNode) Run() {
	go w.ReceiveMultiAnnouncments()
	go w.ListenAndServe()
}

func (w *webSnakeNormalNode) join(msg pubsub.Message) {
	log.Println("Start to join into the game ------------------")
	defer log.Println("------------------ sent join message")
	announce := msg.Msg.GetAnnouncement()

	log.Println("joiner catch command to join and hanlding it")
	log.Println("Try to connect to the game:", announce)

	tp := websnake.PlayerType_HUMAN
	r := websnake.NodeRole_NORMAL

	message := websnake.GameMessage{
		MsgSeq:     new(int64),
		SenderId:   new(int32),
		ReceiverId: new(int32),
		Type: &websnake.GameMessage_Join{
			Join: &websnake.GameMessage_JoinMsg{
				PlayerType:    &tp,
				PlayerName:    &w.username,
				GameName:      announce.GetGames()[0].GameName,
				RequestedRole: &r,
			},
		},
	}

	date, err := proto.Marshal(&message)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	w.SendTo(date, *msg.To)
}

func (mcc *mcastConn) Close() {
	mcc.baseconn.Close()
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
		baseconn:  conn,
		multiconn: p,
		iface:     *netIface,
	}, nil
}

// Run in goroutine.
func (w *webSnakeNormalNode) ReceiveMultiAnnouncments() {

	for {
		buf := make([]byte, 1024)
		n, _, srcAddr, err := w.mcastconn.multiconn.ReadFrom(buf)
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
		log.Println(msg.GetType())
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

func (w *webSnakeNormalNode) SendTo(data []byte, to net.UDPAddr) {
	log.Println("Send Unicast message to ", to.String())
	defer log.Println("---------------------------------------Unicast message sent.")

	n, err1 := w.unicastconn.WriteTo(data, &to)

	if err1 != nil {
		log.Println("Catch err: ", err1)
		os.Exit(1)
	}
	log.Println("Sent:", n)
}

func (w *webSnakeNormalNode) ListenAndServe() {
	log.Println("Master udp conn addr:", w.unicastconn.LocalAddr().String())

	for {
		buf := make([]byte, 4096)
		n, addr, err := w.unicastconn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Err:", err)
			continue
		}
		if n == 0 {
			continue
		}

		buf = buf[0:n]
		message := websnake.GameMessage{}
		if err := proto.Unmarshal(buf, &message); err != nil {
			log.Fatal("Protoerror in ListenAndServe: ", err)
		}

		log.Println("Read from udp socket:", message)

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
		case message.GetPing() != nil:
			{
				// HandlePing(eventMessage)
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

func (w *webSnakeNormalNode) handleAck(eventMessage pubsub.Message) {
	log.Println("Read from udp socket:", eventMessage)
	pubsub.GetGlobalPubSubService().Publish("ack", eventMessage)
}

func (w *webSnakeNormalNode) handleState(eventMessage pubsub.Message) {
	log.Println("Read from udp socket:", eventMessage)
	players := eventMessage.Msg.GetState().State.Players.Players
	addMasterAddresIntoPlayers(players, eventMessage.From)
	pubsub.GetGlobalPubSubService().Publish("newgamestate", eventMessage)
}

func (w *webSnakeNormalNode) sendSteer(msg pubsub.Message) {
	log.Println("Start to send steer $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
	defer log.Println("------------------ sent steer message")

	data, err := proto.Marshal(msg.Msg)
	if err != nil {
		log.Fatal(err)
	}

	w.SendTo(data, *msg.To)
}
