package server

import (
	"main/connection"
	"main/fs"
	"os"
	"path/filepath"
)

type FileConsumer struct {
	filename string
}

func (fc *FileConsumer) HandleFileMetadata(f connection.File) {
	fc.filename = filepath.Base(f.Path)
}

func (fc *FileConsumer) HandleBytes(data []byte) error {
	fs.CreateDirectoryIfNotExists(defaultUploadsDir)

	fc.filename = filepath.Join(defaultUploadsDir, fc.filename)

	outputFilepath := fs.SelectUniqueNameForFile(fc.filename)
	file, err := os.Create(outputFilepath)
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
