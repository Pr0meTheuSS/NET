package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"golang.org/x/net/ipv4"
)

const (
	multicastGroupIPv4 = "224.0.0.1"
	port               = 9876
)

// TODO: run on windows
// TODO unix time Update with system time

var clones = map[string]int64{}

func selectNetInterfaceCli() *net.Interface {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Get net interfaces error:", err)
		os.Exit(1)
	}

	fmt.Println("Net interfaces:")
	for i, iface := range interfaces {
		fmt.Printf("%d: %s\n", i+1, iface.Name)
	}

	fmt.Print("Choose net interface: \n")
	var selectedIndex int
	_, err = fmt.Scan(&selectedIndex)
	if err != nil || selectedIndex < 1 || selectedIndex > len(interfaces) {
		fmt.Println("Cannot choose this interface number.")
		fmt.Println(selectedIndex)
		os.Exit(1)
	}

	selectedInterface := interfaces[selectedIndex-1]
	fmt.Printf("Set net interface: %s\n", selectedInterface.Name)

	return &selectedInterface
}

func main() {
	multicastGroup := multicastGroupIPv4

	if len(os.Args) > 1 {
		multicastGroup = os.Args[1]
	}

	conn, err := net.ListenPacket("udp", multicastGroup+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	netInterface := selectNetInterfaceCli()

	p := ipv4.NewPacketConn(conn)
	p.SetMulticastInterface(netInterface)
	defer p.Close()

	group := net.ParseIP(multicastGroup)
	if group == nil {
		fmt.Println("Invalid multicast group address.")
		os.Exit(1)
	}

	if err := p.JoinGroup(netInterface, &net.UDPAddr{IP: group, Port: port}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Listening for multicast messages on %s:%d...\n", group.String(), port)

	go receiveMulticastMessages(p)
	sendingMulticastMessages(group)
}

func sendingMulticastMessages(group net.IP) {
	sendConn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: group, Port: port})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer sendConn.Close()

	for {
		senderAddress := sendConn.LocalAddr().String()
		if _, err := sendConn.Write([]byte(senderAddress)); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func receiveMulticastMessages(p *ipv4.PacketConn) {
	buf := make([]byte, 1024)

	for {
		oldClonesAmount := len(clones)
		_, _, srcAddr, err := p.ReadFrom(buf)
		if err != nil {
			fmt.Println(err)
			continue
		}

		clones[srcAddr.String()] = int64(time.Now().Unix())

		if oldClonesAmount != len(clones) {
			printClones()
		}

		if cleanup() {
			printClones()
		}
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printClones() {
	clearScreen()
	for k, _ := range clones {
		fmt.Println("Clone address " + k + ";")
	}
}

func cleanup() bool {
	wasSomethingDeleted := false
	nowUnix := time.Now().Unix()
	for k, v := range clones {
		if nowUnix-v > 5 {
			delete(clones, k)
			wasSomethingDeleted = true
		}
	}

	return wasSomethingDeleted
}
