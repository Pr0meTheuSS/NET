package server

import (
	"main/connection"
	"main/fs"
	"os"
	"path/filepath"
)

// FileConsumer implements DataConsumer interface.
type FileConsumer struct {
	filename         string
	alreadyReadBytes int64
}

func (fc *FileConsumer) HandleFileMetadata(f connection.File) string {
	fs.CreateDirectoryIfNotExists(defaultUploadsDir)
	fc.filename = filepath.Base(f.Path)

	fc.filename = filepath.Join(defaultUploadsDir, fc.filename)
	fc.filename = fs.SelectUniqueNameForFile(fc.filename)

	return filepath.Base(fc.filename)
}

func (fc *FileConsumer) HandleBytes(data []byte) error {
	file, err := os.OpenFile(fc.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
