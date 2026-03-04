package main

import (
	"cli-chat-client/bubble"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	f, _ := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)

	conn, err := connect()
	if err != nil {
		log.Fatal("Could not connect to server")
	}

	// reader := bufio.NewReader(os.Stdin)

	p := tea.NewProgram(bubble.InitialModel(
		func(msg string) error {
			log.Println("Sending!")
			_, err := conn.Write([]byte(msg))
			if err != nil {
				return err
			}
			return nil
		},
	))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}

	messages := make(chan []byte)

	go receive(conn, messages)

	go func() {
		for message := range messages {
			log.Println("IncomingMessage", string(message))
			p.Send(bubble.IncomingMessage{
				Sender:  "foo",
				Message: string(message),
			})
		}
	}()

	//
	// for {
	// 	fmt.Print("> ")
	// 	m, err := reader.ReadString('\n')
	// 	if err != nil {
	//
	// 	}
	// 	if len(m) > 0 {
	// 		conn.Write([]byte(m))
	// 	}
	//
	// }

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
