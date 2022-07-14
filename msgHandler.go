package main

import (
	"log"

	"github.com/gorilla/websocket"
)

func msgHandleRoutine() {
	for {
		msg := <-broadcastChannel
		broadcastMessage(msg)
	}
}

func broadcastMessage(msg ChatMessage) {
	for client := range clientMap {
		sendMessage(client, msg)
	}
}

func sendMessage(client *websocket.Conn, msg ChatMessage) {
	err := client.WriteJSON(msg)

	if err != nil {
		log.Printf("s -> c error: %v", err)
		client.Close()
		delete(clientMap, client)
	}
}