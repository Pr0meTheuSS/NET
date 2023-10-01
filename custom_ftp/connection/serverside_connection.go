package connection

import (
	"errors"
	"main/ftp"
	"net"
)

type ServerSideConnection interface {
	receiveHandshake() (*ftp.Handshake, error)
	sendHandshake(*ftp.Handshake) error

	receiveChunk() (*ftp.FileChunk, error)
	handleChunk([]byte) error

	sendResponse(*ftp.FileTransferResponse) error

	isLoaded() bool
	close()
	ServerServe() error
}

type DataConsumer interface {
	HandleBytes([]byte) error
	HandleFileMetadata(File)
}

func NewServerSideConnection(conn net.Conn, dc DataConsumer) ServerSideConnection {
	return &ServerSideConnectionImpl{
		conn:     NewConnection(conn),
		consumer: dc,
	}
}

type ServerSideConnectionImpl struct {
	conn             *Connection
	consumer         DataConsumer
	alreadyReadBytes uint64
}

func (ssc *ServerSideConnectionImpl) close() {
	ssc.conn.Close()
}

func (ssc *ServerSideConnectionImpl) isLoaded() bool {
	return ssc.conn.isLoaded()
}

func (ssc *ServerSideConnectionImpl) sendHandshake(handshake *ftp.Handshake) error {
	return ssc.conn.sendHandshake(handshake)
}

func (ssc *ServerSideConnectionImpl) receiveHandshake() (*ftp.Handshake, error) {
	return ssc.conn.receiveHandshake()
}

func (ssc *ServerSideConnectionImpl) receiveChunk() (*ftp.FileChunk, error) {
	return ssc.conn.receiveChunk()
}

func (ssc *ServerSideConnectionImpl) sendResponse(resp *ftp.FileTransferResponse) error {
	return ssc.conn.sendResponse(resp)
}

func (ssc *ServerSideConnectionImpl) handleChunk(data []byte) error {
	return ssc.consumer.HandleBytes(data)
}

func (ssc *ServerSideConnectionImpl) ServerServe() error {
	defer ssc.close()

	handshake, err := ssc.receiveHandshake()
	if err != nil {
		return err
	}

	ssc.conn.fileSizeBytes = handshake.TotalSize
	ssc.consumer.HandleFileMetadata(File{
		Path:      handshake.Filename,
		SizeBytes: handshake.TotalSize,
	})

	if err := ssc.sendHandshake(handshake); err != nil {
		return err
	}

	if err := ssc.receiveChunks(); err != nil {
		return err
	}

	return ssc.sendResponse(&ftp.FileTransferResponse{Success: true})
}

func (ssc *ServerSideConnectionImpl) receiveChunks() error {
	for chunk, err := ssc.receiveChunk(); ; chunk, err = ssc.receiveChunk() {
		if err != nil {
			if err1 := ssc.sendResponse(&ftp.FileTransferResponse{Success: false}); err1 != nil {
				return errors.Join(err, err1)
			}

			return err
		}

		ssc.handleChunk(chunk.GetData())
		if ssc.isLoaded() {
			break
		}
	}

	return nil
}
