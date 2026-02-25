package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := connect()
	if err != nil {
		log.Fatal("Could not connect to server")
	}
	messages := make(chan []byte)
	go receive(conn, messages)

	go func() {
		for message := range messages {
			fmt.Print("<", string(message))
		}
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		m, err := reader.ReadString('\n')
		if err != nil {

		}
		if len(m) > 0 {
			conn.Write([]byte(m))
		}

	}

}

func connect() (net.Conn, error) {
	return net.Dial("tcp", ":42069")
}

func receive(r io.ReadCloser, messages chan []byte) {

	for {
		buf := make([]byte, 8)
		n, err := r.Read(buf)
		if err != nil {
			fmt.Println()
			log.Fatal("Failed to read with ", err)
		}

		messages <- buf[:n]
	}
}
