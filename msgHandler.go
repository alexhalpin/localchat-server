package main

import (
	"errors"
	"log"
	"time"
)

// func msgHandleRoutine() {
// 	for {
// 		msg := <-broadcastChannel
// 		broadcastMessage(msg)
// 	}
// }

// outgoing message types
type OutgoingChatMessage struct{
	Kind string `json:"kind"`
	Username string `json:"username"`
	Color string `json:"color"`
	Message string `json:"message"`
}

type OutgoingPackage struct {
	SenderUUID string
	Message OutgoingChatMessage
}

type LoginMessage struct {
	Kind string `json:"kind"`
	UUID string `json:"uuid"`
}

// wait 5 seconds for userupdate message from client 
// if userupdate message received return userdata, true
// if no message or incorrect message received return nil, false
func loginCheck(client Client) (*map[string]string, error) {
	conn := client.Socket
	id := client.UUID

	// Send "login" type message as initial handshake
	ogMsg := LoginMessage{
		Kind: "login",
		UUID: id,
	}
	
	if err := conn.WriteJSON(ogMsg); err != nil {
		return nil, err
	}
	log.Printf("Initiated Login: %v", id)

	// Wait 30s max for user to send back a UserUpdateMessage
	var icMsg IncomingMessage
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	
	if err := conn.ReadJSON(&icMsg); err != nil{
		return nil, err
	}

	if icMsg.Kind != "userUpdate" {
		return nil, errors.New("did not receive message of kind 'userUpdate' as login handshake response")
	}

	log.Printf("Successful Login: %v: %v", icMsg.Data["Username"], icMsg.UUID)
	return &icMsg.Data, nil
}

var broadcastChannel = make(chan OutgoingPackage)

func broadcastMessage(pkg OutgoingPackage) {
	for clientUUID := range clientMap {
		
		dist := 0.0
		var err error

		if clientUUID != pkg.SenderUUID {
			dist, err = uuidDist(pkg.SenderUUID, clientUUID)
		}
		
		if err != nil{
			log.Printf("broadcast error: %v", err)
		} else if dist < 1000{
			sendMessage(clientUUID, pkg.Message)
		}
		
	}
}

func sendMessage(clientUUID string, msg OutgoingChatMessage) {
	socket := clientMap[clientUUID].Socket
	err := socket.WriteJSON(msg)
	if err != nil {
		log.Printf("error sending chat: %v", err)
	}
}


