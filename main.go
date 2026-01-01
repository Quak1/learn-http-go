package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const filename = "messages.txt"

func main() {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("couldn't open file")
	}
	defer f.Close()

	linesCh := getLinesChannel(f)
	for line := range linesCh {
		fmt.Println("read:", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)

		line := ""
		for {
			buffer := make([]byte, 8)
			_, err := f.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}

				fmt.Printf("error")
				return
			}

			parts := strings.Split(string(buffer), "\n")

			for _, part := range parts[:len(parts)-1] {
				line += part
				ch <- line
				line = ""
			}

			line += parts[len(parts)-1]
		}
	}()

	return ch
}
