package client

import (
	"io"
	"main/connection"
	"os"
)

// FileDataProvider implements DataProvider interface.
type FileDataProvider struct {
	File             connection.File
	alreadyReadBytes uint64
}

func (d *FileDataProvider) ProvideBytes(size uint32) ([]byte, error) {
	f, err := os.Open(d.File.Path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	f.Seek(int64(d.alreadyReadBytes), 0)

	if d.File.SizeBytes-d.alreadyReadBytes < uint64(size) {
		size = uint32(d.File.SizeBytes - d.alreadyReadBytes)
	}

	data := make([]byte, size)
	n, err := io.ReadFull(f, data)
	if err != nil {
		return nil, err
	}

	d.alreadyReadBytes += uint64(n)
	return data, nil
}
