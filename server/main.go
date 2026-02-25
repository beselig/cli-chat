package main

import (
	"cli-chat-server/clients"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

const addr string = "0.0.0.0:42069"

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Server failed to start with: ", err)
	}

	messagesChannel := make(chan []byte)
	defer close(messagesChannel)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to conenct to %s: %v", addr, err)
		}
		fmt.Println("Client ", conn.RemoteAddr(), "connected")

		client := clients.Connect(conn)
		go readMessages(client, messagesChannel)
		go writeMessages(client, messagesChannel)

	}
}

func readMessages(c clients.Client, channel chan []byte) {
	for {
		buf := make([]byte, 8)
		n, err := c.Conn.Read(buf)

		if errors.Is(err, io.EOF) {
			clients.Disconnect(c.RemoteAddr)
			break
		}
		if err != nil {
			log.Println("Error reading message from ", c.RemoteAddr)
			_, err := fmt.Fprintf(c.Conn, "Error sending your message! Could not send your message! %v", err)
			if err != nil {
				log.Println("Error writing message to: ", c.RemoteAddr)
				break
			}
			break
		}

		channel <- buf[:n]
	}
}

func writeMessages(c clients.Client, messages chan []byte) {
	for {
		select {
		case message := <-messages:
			_, err := c.Conn.Write(message)
			if err != nil {
				fmt.Println("Client ", c.Conn.RemoteAddr(), " disconnected") // TODO: send message to remaining clients
			}
		case <-c.Done:
			c.Conn.Close()
		}

	}
}
