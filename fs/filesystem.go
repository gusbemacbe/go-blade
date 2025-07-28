package fs

import (
	"errors"
	"io"
	"os"
	"time"
)

func Exists(filename string) (bool, error) {
	_, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func LastModified(filename string) time.Time {
	stat, err := os.Stat(filename)

	if err != nil {
		return time.Unix(0, 0)
	}

	return stat.ModTime()
}

func Load(filename string) ([]byte, error) {
	has, err := Exists(filename)

	if err != nil {
		return []byte(""), err
	}

	if !has {
		return []byte(""), errors.New("file does not exists")
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)

	if err != nil {
		return []byte(""), err
	}

	return bytes, nil
}
