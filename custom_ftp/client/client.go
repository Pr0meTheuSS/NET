package client

import (
	"errors"
	"fmt"
	"main/connection"
	"main/fs"
	"net"
)

type Client struct {
	targetIP   string
	targetPort int
}

func NewClient(ip string, port int) *Client {
	return &Client{
		targetIP:   ip,
		targetPort: port,
	}
}

func (c Client) Send(filename string) error {
	f, err := fs.ParseFile(filename)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   net.ParseIP(c.targetIP),
		Port: c.targetPort,
		Zone: "",
	})
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot connect to server. Error:%s", err))
	}

	csc := connection.NewClientSideConnection(conn, &FileDataProvider{File: *f}, *f)
	return csc.ClientServe()
}
