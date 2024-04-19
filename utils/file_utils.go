package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func GetCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

func ReadFileFromPath(path ...string) []byte {
	if len(path) == 0 {
		return nil
	}
	resultPath := filepath.Join(path...)
	log.Info("Read file from path: %s\n", resultPath)
	file, err := os.Open(resultPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	stat, fileStatErr := file.Stat()
	if fileStatErr != nil {
		fmt.Println(fileStatErr)
		return nil
	}
	defer func(file *os.File) {
		fileCloseErr := file.Close()
		if fileCloseErr != nil {
			return
		}
	}(file)

	buffer := make([]byte, stat.Size())
	/*
		for {
			bytesRead, readFileErr := file.Read(buffer)
			if readFileErr != nil {
				if readFileErr != io.EOF {
					fmt.Println(readFileErr)
				}
				break
			}
			fmt.Println(string(buffer[:bytesRead])) // Print content from buffer
		}
	*/
	_, bufioReadErr := bufio.NewReader(file).Read(buffer)
	if bufioReadErr != nil && bufioReadErr != io.EOF {
		fmt.Println(bufioReadErr)
		return nil
	}
	return buffer
}
