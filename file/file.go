package file

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

const (
	go_ext = ".go"
)

func validate(filePath string) (bool, error) {

	if filepath.Ext(filePath) != go_ext {
		return false, errors.New("[INVALID] Given file is not a valid GO file.")
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, errors.New("[INVALID] Given file path does not exist")
		}
		return false, errors.New("[INVALID] Error occurred while accessing the file")
	}

	if info.IsDir() {
		return false, errors.New("[INVALID] Given file path is a directory. File Expected.")
	}

	return true, nil
}

func NewReadCloser(filePath string) (io.ReadCloser, error) {

	ok, err := validate(filePath)
	if !ok {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("[ERROR] Failed to open the file.")
	}

	return file, nil

}
