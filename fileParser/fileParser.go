package fileparser

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// input: a file path
// output: a string of the file contents, if it is a valid .go file

const (
	go_ext = ".go"
)

func ValidateFile(filePath string) (bool, error) {

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

func GOFileReader(filePath string) (io.ReadCloser, error) {

	ok, err := ValidateFile(filePath)
	if !ok {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("[ERROR] Failed to open the file.")
	}

	return file, nil

}
