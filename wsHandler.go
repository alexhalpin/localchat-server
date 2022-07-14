package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clientMap = make(map[*websocket.Conn]bool)

type ChatMessage struct{
	Username string `json:"username"`
	Message string `json:"message"`
}
var broadcastChannel = make(chan ChatMessage)

var upgrader = websocket.Upgrader{}

func wsHandler(w http.ResponseWriter, r *http.Request){

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	log.Println("WS Endpoint Hit")

	conn, err := upgrader.Upgrade(w, r, nil)					// establish websocket connection from initial request at "/ws"
	if err != nil{
		log.Println("Upgrade FAILED: ", err)
		return
	}
	defer func() {
		conn.Close()
		delete(clientMap, conn)
	}()

	clientMap[conn] = true

	for {
		var msg ChatMessage
		err := conn.ReadJSON(&msg)
		if err != nil{
			log.Println("Read FAILED: ", err)
			return
		}
		log.Printf("$ %v: %v", msg.Username, msg.Message)

		broadcastChannel <- msg
	}
}