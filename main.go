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

	line := ""
	for {
		data := make([]byte, 8)
		_, err := f.Read(data)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			fmt.Printf("error")
			break
		}

		parts := strings.Split(string(data), "\n")

		for _, part := range parts[:len(parts)-1] {
			line += part
			fmt.Printf("read: %s\n", line)
			line = ""
		}

		line += parts[len(parts)-1]
	}
}
