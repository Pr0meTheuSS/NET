package websnake

import (
	"net"
	"net/netip"
)

type WebSnakeNode interface {
	SendTo(data []byte, to netip.AddrPort)
	SendMulticast(data []byte, group net.IP)
}
