package main

import (
	"math"
	"strconv"
)

func eucDist(lat1, long1, lat2, long2 float64) float64 {
	R := 6371000.0
	dlong := (long2 - long1) * math.Pi / 180
	dlat := (lat2 - lat1) * math.Pi / 180

	dx := R * math.Cos(lat2) * math.Sin(dlong)
	dy := R * math.Sin(dlat)
	return math.Sqrt(dx*dx + dy*dy)
}

func uuidDist(uuid1 string, uuid2 string) (float64, error) {
	lat1, err := strconv.ParseFloat(clientMap[uuid1].Data.Lat, 64)
	if err != nil {
		return 0.0, err
	}
	long1, err := strconv.ParseFloat(clientMap[uuid1].Data.Long, 64)
	if err != nil {
		return 0.0, err
	}
	lat2, err := strconv.ParseFloat(clientMap[uuid2].Data.Lat, 64)
	if err != nil {
		return 0.0, err
	}
	long2, err := strconv.ParseFloat(clientMap[uuid2].Data.Long, 64)
	if err != nil {
		return 0.0, err
	}

	
	return eucDist(lat1, long1, lat2, long2), nil
}