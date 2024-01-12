package connection

import (
	"errors"
	"main/ftp"
	"net"
)

type ServerSideConnection interface {
	GetAlreadyReadBytes() int64
	ServerServe() error
}

type DataConsumer interface {
	HandleBytes([]byte) error
	HandleFileMetadata(File) string
}

const (
	// 4 Kb
	maxAvailableChunkSize = 4 * 1024
)

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

func (ssc *ServerSideConnectionImpl) GetAlreadyReadBytes() int64 {
	if ssc.isLoaded() {
		return -1
	}

	return int64(ssc.conn.alreadyReadBytes)
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}

	return b
}

func (ssc *ServerSideConnectionImpl) ServerServe() error {
	defer ssc.close()

	handshake, err := ssc.receiveHandshake()
	if err != nil {
		return err
	}

	ssc.conn.fileSizeBytes = handshake.TotalSize
	ssc.conn.chunkSizeBytes = min(maxAvailableChunkSize, handshake.ChunkSize)
	handshake.ChunkSize = ssc.conn.chunkSizeBytes

	handshake.Filename = ssc.consumer.HandleFileMetadata(File{
		Path:      handshake.Filename,
		SizeBytes: handshake.TotalSize,
	})

	if err := ssc.sendHandshake(handshake); err != nil {
		return err
	}

	if err := ssc.receiveChunks(); err != nil {
		if err1 := ssc.sendResponse(&ftp.FileTransferResponse{Success: false}); err1 != nil {
			return errors.Join(err, err1)
		}
		return err
	}

	if ssc.conn.alreadyReadBytes != ssc.conn.fileSizeBytes {
		return ssc.sendResponse(&ftp.FileTransferResponse{Success: false})
	}

	return ssc.sendResponse(&ftp.FileTransferResponse{Success: true})
}

func (ssc *ServerSideConnectionImpl) receiveChunks() error {
	for chunk, err := ssc.receiveChunk(); ; chunk, err = ssc.receiveChunk() {
		if err != nil {
			return err
		}

		if err := ssc.handleChunk(chunk.GetData()); err != nil {
			return err
		}
		if ssc.isLoaded() {
			break
		}
	}

	return nil
}
