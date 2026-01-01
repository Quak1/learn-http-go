package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const addr = "localhost:42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		str, err := r.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		_, err = conn.Write([]byte(str))
		if err != nil {
			log.Fatal(err)
		}
	}
}
