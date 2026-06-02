package main

import (
	"fmt"
	"io"
	"os"

	"github.com/nikkisversion/astral/fileparser"
)

func main() {
	fmt.Println("Hello from Astral! LFG!")

	fileReader, err := fileparser.GOFileReader("testFile.go")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer fileReader.Close()

	if _, err := io.Copy(os.Stdout, fileReader); err != nil {
		fmt.Println("Error:", err)
	}

}
