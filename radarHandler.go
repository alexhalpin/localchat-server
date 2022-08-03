package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type RadarData struct {
	Kind string `json:"kind"`
	TotalActiveUsers int `json:"totalActiveUsers"`
	Locals []Local `json:"locals"`
	
}

type Local struct {
	Username string `json:"username"`
	Color string `json:"color"`
	Coords []float64 `json:"coords"`
	LoginTime string `json:"loginTime"`
}

type RadarPackage struct {
	TargetSocket *websocket.Conn
	Data RadarData
}

var radarChannel = make(chan RadarPackage)

func sendRadarPackage(radarPackage RadarPackage) {
	target := radarPackage.TargetSocket
	radarData := radarPackage.Data

	if err := target.WriteJSON(radarData); err!=nil{
		log.Printf("error sending radar update: %v", err)
	}
}

func radarTickerRoutine() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		loop:
		for {
		select {
			case <- ticker.C:
				radarDetector()
			case <- quit:
				ticker.Stop()
				break loop
			}
		}
	}()
}

const MAXRADIUS float64 = 1000.0

func radarDetector() {
	for targetUUID := range clientMap {
		var locals []Local = []Local{}
		tlats := clientMap[targetUUID].Data.Lat
		tlongs := clientMap[targetUUID].Data.Long
		
		tlat, err := strconv.ParseFloat(tlats, 64)
		if err!=nil{
			log.Printf("Parsefloat Error: %v", err);
		}
		tlong, err := strconv.ParseFloat(tlongs, 64)
		if err!=nil{
			log.Printf("Parsefloat Error: %v", err);
		}

		for otherUUID := range clientMap {
			if targetUUID == otherUUID{
				continue
			}
			olats := clientMap[otherUUID].Data.Lat
			olongs := clientMap[otherUUID].Data.Long

			olat, err := strconv.ParseFloat(olats, 64)
			if err!=nil{
				log.Printf("Parsefloat Error: %v", err);
			}
			olong, err := strconv.ParseFloat(olongs, 64)
			if err!=nil{
				log.Printf("Parsefloat Error: %v", err);
			}
			// log.Printf("%v\t%v\t%v\t%v", tlat, tlong, olat, olong);
			if eucDist(tlat, tlong, olat, olong) < MAXRADIUS {
				username := clientMap[otherUUID].Data.Username
				color := clientMap[otherUUID].Data.Color
				logintime := clientMap[otherUUID].Data.LoginTime
				coords := []float64{olat, olong}
				locals = append(locals, Local{Username: username, Color: color, Coords: coords, LoginTime: logintime})
			}
		}
		
		var msg RadarPackage
		msg.TargetSocket = clientMap[targetUUID].Socket
		msg.Data.Kind = "radarData"
		msg.Data.TotalActiveUsers = len(clientMap)
		msg.Data.Locals = locals

		radarChannel <- msg
	}
}
