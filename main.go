package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	go outboxRoutine()
	go radarTickerRoutine()

	http.HandleFunc("/ws", wsHandler)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
} 