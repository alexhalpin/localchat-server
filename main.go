package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func wsHandler(w http.ResponseWriter, r *http.Request){

	log.Println("WS Endpoint Hit")

	conn, err := upgrader.Upgrade(w, r, nil)					// establish websocket connection from initial request at "/ws"
	if err != nil{
		log.Println("Upgrade FAILED: ", err)
		return
	}
	defer conn.Close()

	uid := uuid.NewString()

	succConnMsg := fmt.Sprintf("Successful WS Connection: %s", uid)
	log.Print(succConnMsg)

	if err := conn.WriteMessage(1, []byte(succConnMsg)); err != nil{
		log.Println(err)
		return
	}

	
	for {
		messageType, message, err := conn.ReadMessage() 
		if err != nil{
			log.Println("Read FAILED: ", err)
			return
		}
		msgString := string(message)
		log.Printf("Received: '%s'", msgString)

		echoResposeStr := fmt.Sprintf("Echo: '%s'", msgString)
		if err := conn.WriteMessage(messageType, []byte(echoResposeStr)); err != nil || message == nil{
			log.Println(err)
			return
		}
		
		log.Printf("Echoed: '%s'", msgString)
	}
}

func main() {

	http.HandleFunc("/ws", wsHandler)

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
	
} 