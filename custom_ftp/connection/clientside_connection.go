package connection

import (
	"errors"
	"fmt"
	"main/ftp"
	"net"
)

type ClientSideConnection interface {
	ClientServe() error
}

type DataProvider interface {
	ProvideBytes(uint32) ([]byte, error)
}

func NewClientSideConnection(conn net.Conn, dp DataProvider, file File) ClientSideConnection {
	return &ClientSideConnectionImpl{
		conn:     NewConnection(conn),
		provider: dp,
		file:     file,
	}
}

type ClientSideConnectionImpl struct {
	conn     *Connection
	provider DataProvider
	file     File
}

func (csc *ClientSideConnectionImpl) close() {
	csc.conn.Close()
}

func (csc *ClientSideConnectionImpl) receiveHandshake() (*ftp.Handshake, error) {
	return csc.conn.receiveHandshake()
}

func (csc *ClientSideConnectionImpl) sendChunk(chunk *ftp.FileChunk) error {
	return csc.sendChunk(chunk)
}

func (csc *ClientSideConnectionImpl) sendHandshake(handshake *ftp.Handshake) error {
	csc.conn.fileSizeBytes = uint64(csc.file.SizeBytes)
	return csc.conn.sendHandshake(handshake)
}

func (csc *ClientSideConnectionImpl) ClientServe() error {
	defer csc.close()

	handshake := ftp.Handshake{
		ChunkSize:      defaultChunkSize,
		TotalSize:      csc.file.SizeBytes,
		FilenameLenfth: uint32(len(csc.file.Path)),
		Filename:       csc.file.Path,
	}

	if err := csc.sendHandshake(&handshake); err != nil {
		return err
	}

	handshakeFromServer, err := csc.receiveHandshake()
	if err != nil {
		return err
	}

	fmt.Printf("Handshake from server:%+v\n", handshakeFromServer)

	csc.conn.chunkSizeBytes = defaultChunkSize

	if err := csc.sendChunks(); err != nil {
		return err
	}

	resp, err := csc.conn.receiveResponse()
	if err != nil {
		return err
	}

	if !resp.Success {
		return errors.New("Transfer failed.")
	}

	return nil
}

func (csc *ClientSideConnectionImpl) sendChunks() error {
	for data, err := csc.provider.ProvideBytes(uint32(csc.conn.chunkSizeBytes)); len(data) != 0; {
		if err != nil {
			return err
		}

		if err := csc.conn.sendChunk(data); err != nil {
			return err
		}
		data, err = csc.provider.ProvideBytes(uint32(csc.conn.chunkSizeBytes))
	}

	return nil
}
