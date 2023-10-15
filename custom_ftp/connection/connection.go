package connection

import (
	"encoding/binary"
	"io"
	"main/ftp"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

type Connection struct {
	connection       net.Conn
	timeoutSeconds   uint
	alreadyReadBytes uint64
	chunkSizeBytes   uint32
	fileSizeBytes    uint64
}

const (
	// 5 seconds
	defaultTimeoutInSeconds = 5
	// 1 Kb
	defaultChunkSize = 1 * 1024
)

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		timeoutSeconds:   defaultTimeoutInSeconds,
		connection:       conn,
		alreadyReadBytes: 0,
	}
}

func (c *Connection) Close() {
	c.connection.Close()
}

func (c *Connection) Read(data []byte) (int, error) {
	c.connection.SetReadDeadline(time.Now().Add(time.Second * time.Duration(c.timeoutSeconds)))
	return io.ReadFull(c.connection, data)
}

func (c *Connection) receiveHandshake() (*ftp.Handshake, error) {
	messageSizeBytes := make([]byte, 4)
	_, err := c.Read(messageSizeBytes)
	if err != nil {
		return nil, err
	}

	messageSize := binary.BigEndian.Uint32(messageSizeBytes)

	message := make([]byte, messageSize)
	_, err = c.Read(message)
	if err != nil {
		return nil, err
	}

	handshake := ftp.Handshake{}
	if err := proto.Unmarshal(message, &handshake); err != nil {
		return nil, err
	}

	return &handshake, nil
}

func (c *Connection) sendHandshake(handshake *ftp.Handshake) error {
	msg, err := proto.Marshal(handshake)
	if err != nil {
		return err
	}

	messageSizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(messageSizeBytes, uint32(len(msg)))
	_, err = c.Write(messageSizeBytes)
	if err != nil {
		return err
	}

	if _, err := c.Write(msg); err != nil {
		return err
	}

	return nil
}

func (c *Connection) receiveChunk() (*ftp.FileChunk, error) {
	messageSizeBytes := make([]byte, 4)
	_, err := c.Read(messageSizeBytes)
	if err != nil {
		return nil, err
	}

	messageSize := binary.BigEndian.Uint32(messageSizeBytes)

	message := make([]byte, messageSize)
	_, err = c.Read(message)
	if err != nil {
		return nil, err
	}

	chunk := ftp.FileChunk{}
	if err := proto.Unmarshal(message, &chunk); err != nil {
		return nil, err
	}

	c.alreadyReadBytes += uint64(len(chunk.GetData()))
	return &chunk, nil
}

func (c *Connection) sendChunk(data []byte) error {
	chunk := ftp.FileChunk{
		Data: data,
	}

	msg, err := proto.Marshal(&chunk)
	if err != nil {
		return err
	}

	messageSizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(messageSizeBytes, uint32(len(msg)))
	_, err = c.Write(messageSizeBytes)
	if err != nil {
		return err
	}

	if _, err := c.Write(msg); err != nil {
		return err
	}

	return nil
}

func (c *Connection) Write(data []byte) (int, error) {
	return c.connection.Write(data)
}

func (c *Connection) isLoaded() bool {
	return c.alreadyReadBytes >= c.fileSizeBytes && c.alreadyReadBytes != 0
}

func (c *Connection) sendResponse(resp *ftp.FileTransferResponse) error {
	msg, err := proto.Marshal(resp)
	if err != nil {
		return err
	}

	messageSizeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(messageSizeBytes, uint32(len(msg)))
	_, err = c.Write(messageSizeBytes)
	if err != nil {
		return err
	}

	if _, err := c.Write(msg); err != nil {
		return err
	}

	return nil
}

func (c *Connection) receiveResponse() (*ftp.FileTransferResponse, error) {
	messageSizeBytes := make([]byte, 4)
	_, err := c.Read(messageSizeBytes)
	if err != nil {
		return nil, err
	}

	messageSize := binary.BigEndian.Uint32(messageSizeBytes)

	message := make([]byte, messageSize)
	_, err = c.Read(message)
	if err != nil {
		return nil, err
	}

	resp := ftp.FileTransferResponse{}
	if err := proto.Unmarshal(message, &resp); err != nil {
	}

	return &resp, nil
}
