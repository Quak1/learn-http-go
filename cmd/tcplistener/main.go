package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const filename = "messages.txt"
const port = ":42069"

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("A connection has been accepted")

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Println(line)
		}

		fmt.Println("-----------Connection has been closed----------")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		defer f.Close()

		line := ""
		for {
			buffer := make([]byte, 8)
			_, err := f.Read(buffer)
			if err != nil {
				if line != "" {
					ch <- line
				}

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
