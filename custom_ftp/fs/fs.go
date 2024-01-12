package fs

import (
	"main/connection"
	"os"

	"github.com/google/uuid"
)

func ParseFile(filepath string) (*connection.File, error) {
	stat, err := os.Stat(filepath)
	if err == nil {
		return &connection.File{
			Path:      stat.Name(),
			SizeBytes: uint64(stat.Size()),
		}, nil
	}

	return nil, err
}

func CreateDirectoryIfNotExists(directoryPath string) error {
	_, err := os.Stat(directoryPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		err := os.MkdirAll(directoryPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func SelectUniqueNameForFile(filepath string) string {
	fpath := filepath
	for isFileExists(fpath) {
		fpath = filepath + generateRandomHash()
	}

	return fpath
}

func isFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func generateRandomHash() string {
	return uuid.New().String()
}
