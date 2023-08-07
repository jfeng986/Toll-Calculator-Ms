package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"Toll-Calculator/types"

	"github.com/minio/websocket"
)

const (
	sendInterval = time.Second
	wsEndpoint   = "ws://localhost:30001/ws"
)

func genLatLon() (float64, float64) {
	return genCoord(), genCoord()
}

func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func genOBUIDs(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = rand.Intn(99999)
	}
	return ids
}

func main() {
	obuIDs := genOBUIDs(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	// defer conn.Close()
	for {
		for i := 0; i < len(obuIDs); i++ {
			lat, lon := genLatLon()
			data := types.OBUData{
				OBUID: obuIDs[i],
				Lat:   lat,
				Lon:   lon,
			}
			fmt.Printf("%+v\n", data)
			err := conn.WriteJSON(&data)
			if err != nil {
				log.Println(err)
			}
		}

		time.Sleep(sendInterval)

	}
}
